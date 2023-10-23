package third

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"io"
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
	logFilePath := c.LogFilePath
	files, err := os.ReadDir(logFilePath)
	if err != nil {
		return err
	}
	tempFiles := make([]string, 0, len(files))
	defer func() {
		for _, file := range tempFiles {
			_ = os.RemoveAll(file)
		}
	}()
	req := third.UploadLogsReq{}
	for _, file := range files {
		if !checkLogPath(file.Name()) {
			continue
		}
		logName := filepath.Join(logFilePath, file.Name())
		filename := fmt.Sprintf("%s.temp_upload.%d", logName, time.Now().UnixMilli())
		tempFiles = append(tempFiles, filename)
		if err := c.fileCopy(logName, filename); err != nil {
			return err
		}
		resp, err := c.fileUploader.UploadFile(ctx, &uploadfile.UploadFileReq{Filepath: filename, Name: file.Name(), Cause: "upload_logs"}, nil)
		if err != nil {
			return err
		}
		req.FileURLs = append(req.FileURLs, &third.FileURL{
			Filename: filename,
			URL:      resp.URL,
		})
	}
	if len(req.FileURLs) == 0 {
		return errs.ErrData.Wrap("not found log file")
	}
	_, err = util.CallApi[third.UploadLogsResp](ctx, constant.UploadLogsRouter, &req)

	return err
}

func checkLogPath(logPath string) bool {
	if len(logPath) < len("open-im-sdk-core.yyyy-mm-dd") {
		return false
	}
	logTime := logPath[len(logPath)-len(".yyyy-mm-dd"):]
	if _, err := time.Parse(".2006-01-02", logTime); err != nil {
		return false
	}
	if !strings.HasPrefix(logPath, "open-im-sdk-core.") {
		return false
	}

	return true
}

func (c *Third) fileCopy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}
