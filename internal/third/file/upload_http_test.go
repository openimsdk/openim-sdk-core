//go:build !(js && wasm)

package file

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/errs"
)

func TestMultipartResumeSkipsStoredPartWithoutHashVerification(t *testing.T) {
	f := newMultipartFixture(t, 8192, 3*8192+17)
	f.seedResume(0)
	cb := &resumeMutationCallback{multipartCallback: multipartCallback{t: t, stored: 8192, parts: make(map[int]int)}, filepath: f.request.Filepath}
	for range len(f.sizes) {
		f.release <- struct{}{}
	}
	if _, err := f.uploader.UploadFile(f.ctx, f.request, cb); err != nil {
		t.Fatal(err)
	}
	if f.puts[1] != 0 || cb.parts[0] != 1 || f.completed != 1 {
		t.Fatal("stored part was not skipped and reported before completion")
	}
}

type resumeMutationCallback struct {
	multipartCallback
	filepath string
}

func (c *resumeMutationCallback) UploadID(string) {
	f, err := os.OpenFile(c.filepath, os.O_WRONLY, 0)
	if err != nil {
		c.t.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			c.t.Error(err)
		}
	}()
	if _, err := f.WriteAt([]byte{255}, 0); err != nil {
		c.t.Fatal(err)
	}
}

func TestMultipartProgressAllowsReorderedAggregateReports(t *testing.T) {
	cb := &multipartCallback{t: t, stored: 10}
	cb.UploadComplete(100, 100, 10)
	cb.UploadComplete(100, 50, 10)
	if cb.progress != 100 {
		t.Fatalf("maximum aggregate progress = %d, want 100", cb.progress)
	}
}

func TestMultipartHTTPRejectsPUTFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "part rejected", http.StatusInternalServerError)
	}))
	defer server.Close()
	request, err := http.NewRequest(http.MethodPut, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = NewFile().doPut(request.Context(), server.Client(), request.URL, nil, strings.NewReader("part"), 4)
	if err == nil || !strings.Contains(err.Error(), "status code 500, body part rejected") {
		t.Fatalf("PUT failure lost: %v", err)
	}
}

func TestMultipartHTTP(t *testing.T) {
	for _, mode := range []string{"concurrent", "resume", "cancel", "failure", "serial"} {
		t.Run(mode, func(t *testing.T) {
			partSize, size := int64(8192), int64(23*8192+17)
			wantOverlap := 4
			if mode == "serial" {
				partSize = 22 * 1024 * 1024
				size = partSize + 17
				wantOverlap = 1
			}
			f := newMultipartFixture(t, partSize, size)
			cb := &multipartCallback{t: t, parts: make(map[int]int)}
			if mode == "resume" {
				f.seedResume(0, 2, 5)
				cb.stored = 3 * partSize
			}
			if mode == "failure" {
				f.failPart = 1
			}
			done := make(chan error, 1)
			go func() {
				resp, err := f.uploader.UploadFile(f.ctx, f.request, cb)
				if err == nil && resp.URL != "https://example.com/multipart" {
					t.Errorf("unexpected completion URL: %s", resp.URL)
				}
				done <- err
			}()
			for range wantOverlap {
				select {
				case <-f.started:
				case <-f.ctx.Done():
					t.Fatalf("PUTs did not overlap: need %d simultaneous requests", wantOverlap)
				}
			}
			if mode == "cancel" {
				if n := f.uploader.CancelUpload(f.ctx, "multipart-operation"); n != 1 {
					t.Fatalf("cancel registrations: %d", n)
				}
			} else {
				for range len(f.sizes) {
					f.release <- struct{}{}
				}
			}
			var err error
			select {
			case err = <-done:
			case <-f.ctx.Done():
				t.Fatal("multipart upload did not return")
			}
			if mode == "cancel" {
				var codeErr errs.CodeError
				if !errors.As(err, &codeErr) || codeErr.Code() != sdkerrs.CancelUploadingFile {
					t.Fatalf("cancel cause lost: %v", err)
				}
			} else if mode == "failure" {
				if err == nil {
					t.Fatal("rejected PUT succeeded")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if cb.progress != size {
					t.Fatalf("final progress %d, want %d", cb.progress, size)
				}
				wantType := 1
				if mode == "resume" {
					wantType = 2
				}
				if cb.typ != wantType {
					t.Fatalf("completion type %d, want %d", cb.typ, wantType)
				}
				for i := range f.sizes {
					if cb.parts[i] != 1 {
						t.Errorf("part %d callbacks: %d", i, cb.parts[i])
					}
				}
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.peak != wantOverlap {
				t.Fatalf("peak PUT concurrency %d, want %d", f.peak, wantOverlap)
			}
			if mode == "cancel" || mode == "failure" {
				if f.completed != 0 || cb.typ != 0 {
					t.Fatal("failed upload completed")
				}
				stored, err := f.database.GetUpload(f.ctx, f.hash)
				if err != nil {
					t.Fatal(err)
				}
				bits, err := base64.StdEncoding.DecodeString(stored.UploadInfo)
				if err != nil {
					t.Fatal(err)
				}
				bitmap := ParseBitmap(bits, len(f.sizes))
				for i := range f.sizes {
					if bitmap.Get(i) {
						t.Errorf("unsuccessful part %d persisted", i)
					}
				}
			} else {
				if f.completed != 1 {
					t.Fatalf("completion calls: %d", f.completed)
				}
				for i := range f.sizes {
					want := 1
					if mode == "resume" && (i == 0 || i == 2 || i == 5) {
						want = 0
					}
					if f.puts[i+1] != want {
						t.Errorf("part %d PUT count: %d, want %d", i+1, f.puts[i+1], want)
					}
				}
				if f.signs < 2 && mode != "serial" {
					t.Fatal("signature refresh not exercised")
				}
			}
			t.Logf("%s: peak overlapping PUTs=%d; ordered MD5s, bytes, progress, SQLite bitmap verified", mode, f.peak)
		})
	}
}
