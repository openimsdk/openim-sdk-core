package third

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
)

func TestUploadLogs(t *testing.T) {
	c := third.NewThird()
	ctx := context.Background()
	params := []sdk_params_callback.UploadLogParams{{LogFilePath: "valid/path"}, {LogFilePath: "invalid/path"}}

	for _, param := range params {
		err := c.UploadLogs(ctx, param)
		if err != nil {
			t.Errorf("UploadLogs() error = %v", err)
		}
	}
}

func TestUploadLogsPrivate(t *testing.T) {
	c := third.NewThird()
	ctx := context.Background()
	params := []sdk_params_callback.UploadLogParams{{LogFilePath: "valid/path"}, {LogFilePath: "invalid/path"}}

	for _, param := range params {
		err := c.uploadLogs(ctx, param)
		if err != nil {
			t.Errorf("uploadLogs() error = %v", err)
		}
	}
}

func TestCheckLogPath(t *testing.T) {
	c := third.NewThird()
	logPaths := []string{"valid/path", "invalid/path"}

	for _, logPath := range logPaths {
		result := c.checkLogPath(logPath)
		if !result {
			t.Errorf("checkLogPath() = %v, want true", result)
		}
	}
}
