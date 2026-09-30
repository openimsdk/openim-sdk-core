package third

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
)

func TestLogTailAndZip(t *testing.T) {
	for _, tc := range []struct {
		name, input, want string
		n                 int
	}{
		// Enterprise counts the terminal newline toward n.
		{"newline", "a\nb\nc\n", "c\n", 2},
		{"no-newline", "a\nb\nc", "b\nc", 2},
		{"single", "a\nb\n", "", 1},
		{"crlf", "a\r\nb\r\nc\r\n", "c\r\n", 2},
		{"short", "a\nb", "a\nb", 9},
		{"empty", "", "", 2},
		{"large-line", "head\n" + strings.Repeat("x", 11*1024*1024) + "\ntail\n", strings.Repeat("x", 11*1024*1024) + "\ntail\n", 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			src, dst := filepath.Join(dir, "source"), filepath.Join(dir, "tail")
			if err := os.WriteFile(src, []byte(tc.input), 0600); err != nil {
				t.Fatal(err)
			}
			if err := writeLastNLines(context.Background(), src, dst, tc.n); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "logs.zip")
			if err := zipFiles(context.Background(), path, []string{dst}); err != nil {
				t.Fatal(err)
			}
			zr, err := zip.OpenReader(path)
			if err != nil {
				t.Fatal(err)
			}
			defer zr.Close()
			r, err := zr.File[0].Open()
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()
			data, err := io.ReadAll(r)
			if err != nil || string(data) != tc.want {
				t.Fatalf("tail differs: size %d, want %d, error %v", len(data), len(tc.want), err)
			}
			ctx, cancel := context.WithCancelCause(context.Background())
			cause := errors.New("stop logs")
			cancel(cause)
			if err := zipFiles(ctx, path, []string{dst}); !errors.Is(err, cause) {
				t.Fatalf("zip cancellation: %v", err)
			}
			wantErr := cause
			if tc.input == "" {
				wantErr = nil // Enterprise returns before reading an empty source.
			}
			if err := writeLastNLines(ctx, src, dst, tc.n); !errors.Is(err, wantErr) {
				t.Fatalf("tail cancellation: %v", err)
			}
		})
	}
}

func TestUploadLogsHTTP(t *testing.T) {
	for _, mode := range []string{"all", "tail", "report-failure", "cancel", "cancel-post"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			name := "open-im-sdk-core.2026-09-30"
			content := "one\ntwo\nthree\n"
			if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			uploaded := make(chan []byte, 1)
			reported := make(chan json.RawMessage, 1)
			postStarted := make(chan struct{}, 1)
			release := make(chan struct{})
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var data string
				switch r.URL.Path {
				case "/object/part_limit":
					data = `{"minPartSize":1048576,"maxPartSize":1048576,"maxNumSize":100}`
				case "/object/initiate_multipart_upload":
					data = fmt.Sprintf(`{"upload":{"uploadID":"logs","partSize":1048576,"sign":{"url":%q,"parts":[{"partNumber":1}]}}}`, server.URL+"/put")
				case "/put":
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
						return
					}
					uploaded <- body
					if mode == "cancel" {
						<-release
					}
					return
				case "/object/complete_multipart_upload":
					if mode == "cancel-post" {
						postStarted <- struct{}{}
						<-release
					}
					data = `{"url":"https://example.com/logs.zip"}`
				case "/third/logs/upload":
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
						return
					}
					reported <- body
					if mode == "report-failure" {
						_, _ = io.WriteString(w, `{"errCode":500,"errMsg":"report failed"}`)
						return
					}
					data = `{}`
				default:
					t.Errorf("unexpected request %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				_, _ = fmt.Fprintf(w, `{"errCode":0,"data":%s}`, data)
			}))
			defer server.Close()
			defer close(release)
			uploader := file.NewFile()
			uploader.SetLoginUserID("user")
			client := NewThird(uploader)
			client.SetLogFilePath(dir)
			client.SetLoginUserID("user")
			client.SetPlatform(5)
			ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{UserID: "user", IMConfig: &sdk_struct.IMConfig{ApiAddr: server.URL}})
			ctx = ccontext.WithOperationID(ctx, "logs-qa")
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			line := 2
			if mode == "all" {
				line = 0
			}
			done := make(chan error, 1)
			go func() { done <- client.UploadLogs(ctx, line, "cancel-logs", "qa", nil) }()
			var body []byte
			select {
			case body = <-uploaded:
			case <-ctx.Done():
				t.Fatal("PUT not reached")
			}
			if mode == "cancel-post" {
				select {
				case <-postStarted:
				case <-ctx.Done():
					t.Fatal("POST not reached")
				}
			}
			canceled := mode == "cancel" || mode == "cancel-post"
			if canceled && client.CancelUpload(ctx, "cancel-logs") == 0 {
				t.Fatal("upload not registered")
			}
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("upload did not return")
			}
			if (err != nil) != (canceled || mode == "report-failure") {
				t.Fatalf("unexpected result: %v", err)
			}
			if canceled {
				var codeErr errs.CodeError
				wantCode := sdkerrs.CancelUploadingFile
				if !errors.As(err, &codeErr) || codeErr.Code() != wantCode {
					t.Fatalf("cancellation code differs from enterprise: want %d, error %v", wantCode, err)
				}
			}
			if !canceled {
				select {
				case report := <-reported:
					if !bytes.Contains(report, []byte("https://example.com/logs.zip")) || !bytes.Contains(report, []byte(`"ex":"qa"`)) {
						t.Fatalf("invalid report: %s", report)
					}
				default:
					t.Fatal("missing log report")
				}
			} else if len(reported) != 0 {
				t.Fatal("canceled upload reported")
			}
			zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				t.Fatal(err)
			}
			if len(zr.File) != 1 {
				t.Fatalf("zip entries: %d", len(zr.File))
			}
			r, err := zr.File[0].Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(r)
			r.Close()
			want := "three\n"
			if mode == "all" {
				want = content
			}
			if err != nil || string(data) != want {
				t.Fatalf("uploaded log: %q, error %v", data, err)
			}
			source, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil || string(source) != content {
				t.Fatal("source log changed")
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 1 {
				t.Fatalf("temporary files remain: %v, %v", entries, err)
			}
			if n := client.CancelUpload(ctx, "cancel-logs"); n != 0 {
				t.Fatalf("stale cancel registrations: %d", n)
			}
			if err := client.UploadLogs(ctx, -1, "", "", nil); err == nil {
				t.Fatal("negative line accepted")
			}
			t.Logf("%s: ZIP PUT verified; report, cancellation and cleanup verified", mode)
		})
	}
}
