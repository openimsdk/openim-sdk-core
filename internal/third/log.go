package third

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

const (
	buffer = 10 * 1024 * 1024
)

func (c *Third) uploadLogs(ctx context.Context, line int, ex string, progress Progress) (err error) {
	if c.logUploadLock.TryLock() {
		defer c.logUploadLock.Unlock()
	} else {
		return errs.New("log file is uploading").Wrap()
	}
	if line < 0 {
		return errs.New("line is illegal").Wrap()
	}

	logFilePath := c.LogFilePath
	entrys, err := os.ReadDir(logFilePath)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entrys))
	switch line {
	case 0:
		// all logs
		for _, entry := range entrys {
			if (!entry.IsDir()) && (!strings.HasSuffix(entry.Name(), ".zip")) && checkLogPath(entry.Name()) {
				files = append(files, filepath.Join(logFilePath, entry.Name()))
			}
		}
		if len(files) == 0 {
			return errs.New("not found log file").Wrap()
		}
		defer func() {
			if err == nil {
				// remove old file
				for _, f := range files[:len(files)-1] {
					if err := os.Remove(f); err != nil {
						log.ZError(ctx, "remove file failed", err, "file name", f)
					}
				}
				// truncate now log file
				f, err := os.OpenFile(files[len(files)-1], os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					log.ZError(ctx, "remove file failed", err, "file name", f)
				}
				_ = f.Close()
			}
		}()
	default:
		for i := len(entrys) - 1; i >= 0; i-- {
			// get newest log file
			if (!entrys[i].IsDir()) && (!strings.HasSuffix(entrys[i].Name(), ".zip")) && checkLogPath(entrys[i].Name()) {
				files = append(files, filepath.Join(logFilePath, entrys[i].Name()))
				break
			}
		}
		if len(files) == 0 {
			return errs.New("not found log file").Wrap()
		}
		lines, err := readLastNLines(files[0], line)
		if err != nil {
			return err
		}
		data := strings.Join(lines, "\n")

		// create tmp file
		filename := fmt.Sprintf("%s_temp%s", strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0])), filepath.Ext(files[0]))
		files[0] = filepath.Join(logFilePath, filename)
		err = os.WriteFile(files[0], []byte(data), 0644)
		if err != nil {
			return errs.Wrap(err)
		}
		defer func() {
			if err := os.Remove(files[0]); err != nil {
				log.ZError(ctx, "remove file failed", err, "file name", files[0])
			}
		}()
	}

	zippath := filepath.Join(logFilePath, fmt.Sprintf("%d_%d.zip", time.Now().UnixMilli(), rand.Uint32()))
	defer os.Remove(zippath)
	if err := zipFiles(zippath, files); err != nil {
		return err
	}
	reqUpload := &file.UploadFileReq{Filepath: zippath, Name: fmt.Sprintf("sdk_log_%s_%s_%s_%s_%s",
		c.loginUserID, c.systemType, constant.PlatformID2Name[int(c.platformID)], version.Version, filepath.Base(zippath)), Cause: "sdklog", ContentType: "application/zip"}
	resp, err := c.fileUploader.UploadFile(ctx, reqUpload, &progressConvert{ctx: ctx, p: progress})
	if err != nil {
		return err
	}
	ccontext.Info(ctx)
	reqLog := &third.UploadLogsReq{
		Platform:     c.platformID,
		AppFramework: c.systemType,
		Version:      version.Version,
		FileURLs:     []*third.FileURL{{Filename: zippath, URL: resp.URL}},
		Ex:           ex,
	}
	return api.UploadLogs.Execute(ctx, reqLog)
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

func readLastNLines(filename string, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, n)
	count := 0

	scanner := bufio.NewScanner(f)
	buf := make([]byte, buffer)
	scanner.Buffer(buf, buffer)

	for scanner.Scan() {
		lines[count%n] = scanner.Text()
		count++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	start := count - n
	if start < 0 {
		start = 0
	}
	result := make([]string, 0, n)
	for i := start; i < count; i++ {
		result = append(result, lines[i%n])
	}

	return result, nil
}

func (c *Third) printLog(ctx context.Context, logLevel int, file string, line int, msg, err string, keysAndValues []any) {
	errString := errs.New(err)

	log.SDKLog(ctx, logLevel, file, line, msg, errString, keysAndValues)
}
