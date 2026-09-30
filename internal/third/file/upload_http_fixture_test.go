//go:build !(js && wasm)

package file

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/third"
)

type multipartFixture struct {
	t                              *testing.T
	uploader                       *File
	database                       *db.DataBase
	ctx                            context.Context
	request                        *UploadFileReq
	partSize                       int64
	sizes                          []int64
	hashes                         []string
	hash                           string
	started                        chan int
	release                        chan struct{}
	mu                             sync.Mutex
	active, peak, signs, completed int
	puts                           map[int]int
	failPart                       int
	passPart                       int
	failAPI                        string
}

func newMultipartFixture(t *testing.T, partSize, size int64) *multipartFixture {
	t.Helper()
	dir := t.TempDir()
	filename := filepath.Join(dir, "multipart.bin")
	source, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	if err := source.Truncate(size); err != nil {
		t.Fatal(err)
	}
	if size < 1024*1024 {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i*31 + i/8192)
		}
		if _, err := source.WriteAt(data, 0); err != nil {
			t.Fatal(err)
		}
	}
	f := &multipartFixture{t: t, uploader: NewFile(), partSize: partSize, started: make(chan int, 100), release: make(chan struct{}, 100), puts: make(map[int]int)}
	f.request = &UploadFileReq{Filepath: filename, Name: "multipart.bin", CancelID: "multipart-operation"}
	for offset := int64(0); offset < size; offset += partSize {
		length := min(partSize, size-offset)
		h := md5.New()
		if _, err := io.Copy(h, io.NewSectionReader(source, offset, length)); err != nil {
			t.Fatal(err)
		}
		f.sizes = append(f.sizes, length)
		f.hashes = append(f.hashes, hex.EncodeToString(h.Sum(nil)))
	}
	f.hash = f.uploader.partMD5(f.hashes)
	server := httptest.NewServer(http.HandlerFunc(f.serveHTTP))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(f.release) })
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{UserID: "multipart-user", IMConfig: &sdk_struct.IMConfig{ApiAddr: server.URL}})
	ctx = ccontext.WithOperationID(ctx, "multipart-operation")
	var cancel context.CancelFunc
	f.ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	t.Cleanup(cancel)
	f.database, err = db.NewDataBase(f.ctx, "multipart-user", dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := f.database.Close(context.WithoutCancel(f.ctx)); err != nil {
			t.Error(err)
		}
	})
	f.uploader.SetDataBase(f.database)
	f.uploader.SetLoginUserID("multipart-user")
	return f
}

func (f *multipartFixture) seedResume(indices ...int) {
	f.t.Helper()
	bitmap := NewBitmap(len(f.sizes))
	for _, index := range indices {
		bitmap.Set(index)
	}
	if err := f.database.InsertUpload(f.ctx, &model_struct.LocalUpload{PartHash: f.hash, UploadID: "multipart", UploadInfo: base64.StdEncoding.EncodeToString(bitmap.Serialize()), ExpireTime: time.Now().Add(2 * time.Hour).UnixMilli()}); err != nil {
		f.t.Fatal(err)
	}
}

func (f *multipartFixture) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == f.failAPI {
		http.Error(w, "API rejected", http.StatusInternalServerError)
		return
	}
	var data any
	switch r.URL.Path {
	case "/object/part_limit":
		data = &third.PartLimitResp{MinPartSize: f.partSize, MaxPartSize: f.partSize, MaxNumSize: 100}
	case "/object/initiate_multipart_upload":
		data = &third.InitiateMultipartUploadResp{Upload: &third.UploadInfo{UploadID: "multipart", PartSize: f.partSize, ExpireTime: time.Now().Add(2 * time.Hour).UnixMilli(), Sign: &third.AuthSignParts{}}}
	case "/object/auth_sign":
		var req third.AuthSignReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			f.t.Error(err)
			return
		}
		parts := make([]*third.SignPart, 0, len(req.PartNumbers))
		for _, number := range req.PartNumbers {
			parts = append(parts, &third.SignPart{PartNumber: number, Url: fmt.Sprintf("http://%s/put/%d", r.Host, number)})
		}
		f.mu.Lock()
		f.signs++
		f.mu.Unlock()
		data = &third.AuthSignResp{Parts: parts}
	case "/object/complete_multipart_upload":
		var req third.CompleteMultipartUploadReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			f.t.Error(err)
			return
		}
		if strings.Join(req.Parts, ",") != strings.Join(f.hashes, ",") {
			f.t.Error("completion MD5 order differs")
		}
		stored, err := f.database.GetUpload(f.ctx, f.hash)
		if err != nil {
			f.t.Error(err)
			return
		}
		bits, err := base64.StdEncoding.DecodeString(stored.UploadInfo)
		if err != nil {
			f.t.Error(err)
			return
		}
		bitmap := ParseBitmap(bits, len(f.sizes))
		for i := range f.sizes {
			if !bitmap.Get(i) {
				f.t.Errorf("part %d not persisted", i)
			}
		}
		f.mu.Lock()
		f.completed++
		f.mu.Unlock()
		data = &third.CompleteMultipartUploadResp{Url: "https://example.com/multipart"}
	default:
		if !strings.HasPrefix(r.URL.Path, "/put/") {
			http.NotFound(w, r)
			return
		}
		number, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/put/"))
		if err != nil || number < 1 || number > len(f.sizes) {
			f.t.Errorf("invalid part: %s", r.URL.Path)
			return
		}
		f.mu.Lock()
		f.active++
		f.peak = max(f.peak, f.active)
		f.puts[number]++
		f.mu.Unlock()
		defer func() { f.mu.Lock(); f.active--; f.mu.Unlock() }()
		h := md5.New()
		size, err := io.Copy(h, r.Body)
		if err != nil && r.Context().Err() != nil {
			return
		}
		if err != nil || size != f.sizes[number-1] || hex.EncodeToString(h.Sum(nil)) != f.hashes[number-1] || r.ContentLength != size {
			f.t.Errorf("part %d bytes/MD5/content length mismatch: size %d, error %v", number, size, err)
		}
		f.started <- number
		if number == f.passPart {
			return
		}
		if f.failPart != 0 && number != f.failPart {
			select {
			case <-f.ctx.Done():
			case <-r.Context().Done():
			}
		} else {
			select {
			case <-f.release:
			case <-f.ctx.Done():
			case <-r.Context().Done():
			}
		}
		if number == f.failPart {
			http.Error(w, "part rejected", http.StatusInternalServerError)
		}
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		f.t.Error(err)
		return
	}
	if _, err := fmt.Fprintf(w, `{"errCode":0,"data":%s}`, payload); err != nil {
		f.t.Error(err)
	}
}

type multipartCallback struct {
	emptyUploadCallback
	progress, stored int64
	parts            map[int]int
	t                *testing.T
	typ              int
}

func (c *multipartCallback) UploadComplete(size, streamed, stored int64) {
	if streamed < stored || streamed > size || stored != c.stored {
		c.t.Errorf("invalid aggregate progress: %d after %d, storage %d", streamed, c.progress, stored)
	}
	c.progress = max(c.progress, streamed)
}

func (c *multipartCallback) UploadPartComplete(index int, size int64, hash string) { c.parts[index]++ }
func (c *multipartCallback) Complete(size int64, url string, typ int)              { c.typ = typ }
