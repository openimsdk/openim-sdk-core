package third

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenIMSDK/protocol/third"
	uploadfile "github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
)

func (c *Third) UploadLogs(ctx context.Context, params []sdk_params_callback.UploadLogParams) error {
	return c.uploadLogs(ctx, params)
}

func (c *Third) uploadLogs(ctx context.Context, params []sdk_params_callback.UploadLogParams) error {
	logFilePath := c.LogFilePath
	files, err := os.ReadDir(logFilePath)
	if err != nil {
		return err
	}
	req := third.UploadLogsReq{}
	for _, file := range files {
		if !checkLogPath(file.Name()) {
			continue
		}
		var filename = filepath.Join(logFilePath, file.Name())
		resp, err := c.fileUploader.UploadFile(ctx, &uploadfile.UploadFileReq{Filepath: filename, Name: file.Name(), Cause: "upload_logs"}, nil)
		if err != nil {
			return err
		}
		var fileURL third.FileURL
		fileURL.Filename = filename
		fileURL.URL = resp.URL
		req.FileURLs = append(req.FileURLs, &fileURL)
	}
	_, err = util.CallApi[third.UploadLogsResp](ctx, constant.UploadLogsRouter, &req)
	return err
}

func checkLogPath(logpath string) bool {
	if len(logpath) < len("open-im-sdk-core.yyyy-mm-dd") {
		return false
	}
	logTime := logpath[len(logpath)-len(".yyyy-mm-dd"):]
	if _, err := time.Parse(".2006-01-02", logTime); err != nil {
		return false
	}
	if !strings.HasPrefix(logpath, "open-im-sdk-core.") {
		return false
	}

	return true
}
