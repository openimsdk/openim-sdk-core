//go:build !(js && wasm)

package file

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/errs"
)

func TestMultipartHTTPMemoryBudget(t *testing.T) {
	for _, tc := range []struct {
		name     string
		partSize int64
		workers  int
	}{
		{"three-workers", 16 * 1024 * 1024, 3},
		{"two-workers", 21 * 1024 * 1024, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newMultipartFixture(t, tc.partSize, 3*tc.partSize+17)
			cb := &multipartCallback{t: t, parts: make(map[int]int)}
			done := make(chan error, 1)
			go func() { _, err := f.uploader.UploadFile(f.ctx, f.request, cb); done <- err }()
			for range tc.workers {
				select {
				case <-f.started:
				case <-f.ctx.Done():
					t.Fatal("expected workers did not overlap")
				}
			}
			for range len(f.sizes) {
				f.release <- struct{}{}
			}
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-f.ctx.Done():
				t.Fatal("upload did not return")
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.peak != tc.workers || f.completed != 1 || cb.progress != 3*tc.partSize+17 {
				t.Fatalf("peak=%d complete=%d progress=%d", f.peak, f.completed, cb.progress)
			}
			t.Logf("part size %d: verified %d overlapping PUTs", tc.partSize, f.peak)
		})
	}
}

func TestMultipartHTTPAPIFailures(t *testing.T) {
	for _, path := range []string{"/object/part_limit", "/object/initiate_multipart_upload", "/object/auth_sign", "/object/complete_multipart_upload"} {
		t.Run(strings.TrimPrefix(path, "/object/"), func(t *testing.T) {
			f := newMultipartFixture(t, 8192, 3*8192+17)
			f.failAPI = path
			for range len(f.sizes) {
				f.release <- struct{}{}
			}
			cb := &multipartCallback{t: t, parts: make(map[int]int)}
			resp, err := f.uploader.UploadFile(f.ctx, f.request, cb)
			if err == nil || resp != nil || cb.typ != 0 {
				t.Fatalf("failed API produced success: response=%v error=%v type=%d", resp, err, cb.typ)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.completed != 0 {
				t.Fatal("failed API completed upload")
			}
			if path != "/object/complete_multipart_upload" && len(f.puts) != 0 {
				t.Fatal("PUT started despite prerequisite API failure")
			}
			if path == "/object/complete_multipart_upload" {
				for i := range f.sizes {
					if f.puts[i+1] != 1 {
						t.Fatalf("part %d PUT count=%d", i+1, f.puts[i+1])
					}
				}
			}
		})
	}
}

type persistedPartCallback struct {
	multipartCallback
	finished chan int
}

func (c *persistedPartCallback) UploadPartComplete(index int, size int64, hash string) {
	c.multipartCallback.UploadPartComplete(index, size, hash)
	c.finished <- index
}

func TestMultipartHTTPCancelThenResume(t *testing.T) {
	f := newMultipartFixture(t, 8192, 7*8192+17)
	f.passPart = 1
	cb := &persistedPartCallback{multipartCallback: multipartCallback{t: t, parts: make(map[int]int)}, finished: make(chan int, 8)}
	done := make(chan error, 1)
	go func() { _, err := f.uploader.UploadFile(f.ctx, f.request, cb); done <- err }()
	select {
	case index := <-cb.finished:
		if index != 0 {
			t.Fatalf("unexpected persisted part %d", index)
		}
	case <-f.ctx.Done():
		t.Fatal("first part did not persist")
	}
	if n := f.uploader.CancelUpload(f.ctx, f.request.CancelID); n != 1 {
		t.Fatalf("cancel registrations=%d", n)
	}
	select {
	case err := <-done:
		var codeErr errs.CodeError
		if !errors.As(err, &codeErr) || codeErr.Code() != sdkerrs.CancelUploadingFile {
			t.Fatalf("cancel cause=%v", err)
		}
	case <-f.ctx.Done():
		t.Fatal("cancel did not stop upload")
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
		if bitmap.Get(i) != (i == 0) {
			t.Fatalf("unexpected persisted bit %d", i)
		}
	}
	for range 2 * len(f.sizes) {
		f.release <- struct{}{}
	}
	resumed := &multipartCallback{t: t, stored: 8192, parts: make(map[int]int)}
	if _, err := f.uploader.UploadFile(f.ctx, f.request, resumed); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.puts[1] != 1 || f.completed != 1 || resumed.typ != 2 || resumed.progress != 7*8192+17 {
		t.Fatalf("resume: stored part PUTs=%d completed=%d type=%d progress=%d", f.puts[1], f.completed, resumed.typ, resumed.progress)
	}
	for i := range f.sizes {
		if resumed.parts[i] != 1 {
			t.Fatalf("resume part %d callback count=%d", i, resumed.parts[i])
		}
	}
}

func TestMultipartHTTPFullyStoredResume(t *testing.T) {
	f := newMultipartFixture(t, 8192, 3*8192+17)
	f.seedResume(0, 1, 2, 3)
	cb := &multipartCallback{t: t, stored: 3*8192 + 17, parts: make(map[int]int)}
	if _, err := f.uploader.UploadFile(f.ctx, f.request, cb); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.puts) != 0 || f.signs != 0 || f.completed != 1 || cb.typ != 2 {
		t.Fatalf("fully stored resume: PUTs=%v signs=%d completed=%d type=%d", f.puts, f.signs, f.completed, cb.typ)
	}
	for i := range f.sizes {
		if cb.parts[i] != 1 {
			t.Fatalf("part %d callback count=%d", i, cb.parts[i])
		}
	}
}
