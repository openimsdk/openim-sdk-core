package third

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (c *Third) uploadLogs(ctx context.Context, line int, cancelID string, ex string, progress Progress) (err error) {
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
		// create tmp file
		filename := fmt.Sprintf("%s_temp%s", strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0])), filepath.Ext(files[0]))
		src := files[0]
		files[0] = filepath.Join(logFilePath, filename)
		if err := writeLastNLines(ctx, src, files[0], line); err != nil {
			return errs.Wrap(err)
		}
		defer func() {
			if err := os.Remove(files[0]); err != nil {
				log.ZError(ctx, "remove file failed", err, "file name", files[0])
			}
		}()
	}

	return c.uploadZipAndReport(ctx, logFilePath, files, "sdk_log", "sdklog", ex, cancelID, progress)
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

func writeLastNLines(ctx context.Context, filename string, outputPath string, n int) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()
	info, err := src.Stat()
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	if n <= 0 || info.Size() == 0 {
		return nil
	}
	start, err := lastNLinesStart(ctx, src, info.Size(), n)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, &ctxReader{
		ctx:    ctx,
		reader: io.NewSectionReader(src, start, info.Size()-start),
	})
	return err
}

func lastNLinesStart(ctx context.Context, f *os.File, size int64, n int) (int64, error) {
	const blockSize = 64 * 1024
	start := int64(0)
	offset := size
	newlines := 0
	buf := make([]byte, blockSize)
	for offset > 0 {
		select {
		case <-ctx.Done():
			return 0, context.Cause(ctx)
		default:
		}
		readSize := int64(blockSize)
		if offset < readSize {
			readSize = offset
		}
		offset -= readSize
		chunk := buf[:readSize]
		if _, err := f.ReadAt(chunk, offset); err != nil && err != io.EOF {
			return 0, err
		}
		for i := len(chunk) - 1; i >= 0; i-- {
			if chunk[i] != '\n' {
				continue
			}
			newlines++
			if newlines == n {
				start = offset + int64(i) + 1
				offset = 0
				break
			}
		}
	}
	return start, nil
}

func (c *Third) printLog(ctx context.Context, logLevel int, file string, line int, msg, err string, keysAndValues []any) {
	errString := errs.New(err)

	log.SDKLog(ctx, logLevel, file, line, msg, errString, keysAndValues)
}
