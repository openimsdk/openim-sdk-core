package third

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
)

func TestUploadLogs(t *testing.T) {
	c := &Third{LogFilePath: "/tmp/logs"}

	// Test when the log file path is valid and contains log files
	os.MkdirAll(c.LogFilePath, 0755)
	ioutil.WriteFile(filepath.Join(c.LogFilePath, "open-im-sdk-core.2022-01-01"), []byte("test"), 0644)
	err := c.UploadLogs(context.Background(), []sdk_params_callback.UploadLogParams{})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// Test when the log file path is valid but does not contain log files
	os.RemoveAll(c.LogFilePath)
	os.MkdirAll(c.LogFilePath, 0755)
	err = c.UploadLogs(context.Background(), []sdk_params_callback.UploadLogParams{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Test when the log file path is invalid
	c.LogFilePath = "/invalid/path"
	err = c.UploadLogs(context.Background(), []sdk_params_callback.UploadLogParams{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCheckLogPath(t *testing.T) {
	c := &Third{}

	// Test when the log file path is valid
	valid := c.checkLogPath("open-im-sdk-core.2022-01-01")
	if !valid {
		t.Errorf("expected true, got false")
	}

	// Test when the log file path is invalid
	valid = c.checkLogPath("invalid.log")
	if valid {
		t.Errorf("expected false, got true")
	}

	// Test when the log file path is empty
	valid = c.checkLogPath("")
	if valid {
		t.Errorf("expected false, got true")
	}
}
