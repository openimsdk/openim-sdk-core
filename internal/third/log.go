package third

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (c *Third) UploadLogs(ctx context.Context, progress Progress) error {
	logFilePath := c.LogFilePath
	entrys, err := os.ReadDir(logFilePath)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entrys))
	for _, entry := range entrys {
		if (!entry.IsDir()) && (!strings.HasSuffix(entry.Name(), ".zip")) && checkLogPath(entry.Name()) {
			files = append(files, filepath.Join(logFilePath, entry.Name()))
		}
	}
	if len(files) == 0 {
		return errors.New("not found log file")
	}
	zippath := filepath.Join(logFilePath, fmt.Sprintf("%d_%d.zip", time.Now().UnixMilli(), rand.Uint32()))
	defer os.Remove(zippath)
	if err := zipFiles(zippath, files); err != nil {
		return err
	}
	reqUpload := &file.UploadFileReq{Filepath: zippath, Name: fmt.Sprintf("sdk_log_%s_%s", c.loginUserID, filepath.Base(zippath)), Cause: "sdklog", ContentType: "application/zip"}
	resp, err := c.fileUploader.UploadFile(ctx, reqUpload, &progressConvert{ctx: ctx, p: progress})
	if err != nil {
		return err
	}
	ccontext.Info(ctx)
	reqLog := &third.UploadLogsReq{
		Platform:   c.platformID,
		SystemType: c.systemType,
		Version:    c.version,
		FileURLs:   []*third.FileURL{{Filename: zippath, URL: resp.URL}},
	}
	_, err = util.CallApi[third.UploadLogsResp](ctx, constant.UploadLogsRouter, reqLog)
	return err
}

// RemoveAllLogFiles delete all log files
func (c *Third) RemoveAllLogFiles(ctx context.Context) error {
	entries, err := os.ReadDir(c.LogFilePath)
	if err != nil {
		return err
	}
	var filePaths []string
	for _, entry := range entries {
		if (!entry.IsDir()) && (!strings.HasSuffix(entry.Name(), ".zip")) && checkLogPath(entry.Name()) {
			filePaths = append(filePaths, filepath.Join(c.LogFilePath, entry.Name()))
		} else {
			log.ZWarn(ctx, "removeLogFiles failed", nil, "entry", entry)
		}
	}
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			log.ZWarn(ctx, "removeLogFiles failed", err, "filePath", filePath)
			continue
		} else {
			log.ZDebug(ctx, "removeLogFiles success", "filePath", filePath)
		}
	}
	return nil
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
	_ = os.RemoveAll(dst)
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
