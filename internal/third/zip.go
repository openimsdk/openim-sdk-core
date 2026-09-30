package third

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zip"
)

func zipFiles(ctx context.Context, outputPath string, files []string) (err error) {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if closeErr := zipWriter.Close(); err == nil {
			err = closeErr
		}
	}()

	for _, file := range files {
		if err := addFileToZip(ctx, zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(ctx context.Context, zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, &ctxReader{
		ctx:    ctx,
		reader: io.LimitReader(file, info.Size()),
	})
	return err
}

type ctxReader struct {
	ctx    context.Context
	reader io.Reader
}

func (c *ctxReader) Read(p []byte) (int, error) {
	select {
	case <-c.ctx.Done():
		return 0, context.Cause(c.ctx)
	default:
	}
	return c.reader.Read(p)
}
