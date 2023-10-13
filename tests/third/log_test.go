package third

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
)

func TestUploadLogs(t *testing.T) {
	c := third.New()
	ctx := context.Background()
	params := []sdk_params_callback.UploadLogParams{
		{
			LogFilePath: "valid/path",
		},
		{
			LogFilePath: "invalid/path",
		},
	}

	for _, param := range params {
		err := c.UploadLogs(ctx, param)
		if err != nil {
			t.Errorf("UploadLogs() error = %v", err)
		}
	}
}

func TestUploadLogsInvalid(t *testing.T) {
	c := third.New()
	ctx := context.Background()
	params := []sdk_params_callback.UploadLogParams{
		{
			LogFilePath: "",
		},
	}

	for _, param := range params {
		err := c.UploadLogs(ctx, param)
		if err == nil {
			t.Errorf("UploadLogs() expected error, got nil")
		}
	}
}

func TestCheckLogPath(t *testing.T) {
	c := third.New()
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "valid path",
			path: "open-im-sdk-core.2022-01-01",
			want: true,
		},
		{
			name: "invalid path",
			path: "invalid/path",
			want: false,
		},
		{
			name: "path not following format",
			path: "open-im-sdk-core.01-01-2022",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.CheckLogPath(tt.path); got != tt.want {
				t.Errorf("CheckLogPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
