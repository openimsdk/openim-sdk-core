package third

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func (c *Third) addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
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
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, io.LimitReader(file, info.Size()))
	return err
}

func zipFiles(zipPath string, files []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	addFileToZip := func(fp string) error {
		file, err := os.Open(fp)
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
		header.Name = filepath.Base(file.Name())
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, io.LimitReader(file, info.Size()))
		return err
	}
	for _, file := range files {
		err := addFileToZip(file)
		if err != nil {
			return err
		}
	}
	if err := zipWriter.Flush(); err != nil {
		return err
	}
	return nil
}
