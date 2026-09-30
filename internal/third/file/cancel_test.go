package file

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/errs"
)

func TestUploadCancelRegistrations(t *testing.T) {
	f := NewFile()
	first, stopFirst := f.WithUploadCancel(context.Background(), "same")
	defer stopFirst()
	second, stopSecond := f.WithUploadCancel(context.Background(), "same")
	stopSecond()
	stopSecond()
	if second.Err() != context.Canceled || first.Err() != nil {
		t.Fatal("cleanup must cancel only its own registration")
	}
	if n := f.CancelUpload(context.Background(), "same"); n != 1 {
		t.Fatalf("canceled %d registrations, want 1", n)
	}
	var codeErr errs.CodeError
	if cause := context.Cause(first); !errors.As(cause, &codeErr) || codeErr.Code() != sdkerrs.CancelUploadingFile {
		t.Fatalf("unexpected cancel cause: %v", cause)
	}
	third, stopThird := f.WithUploadCancel(context.Background(), "same")
	defer stopThird()
	stopFirst()
	if n := f.CancelUpload(context.Background(), "same"); n != 1 || third.Err() == nil {
		t.Fatal("old cleanup removed a new registration")
	}
	if n := f.CancelUpload(context.Background(), "same"); n != 0 {
		t.Fatalf("already canceled: %d", n)
	}
}
