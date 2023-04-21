package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
)

func hashReader(r io.Reader) (string, error) {
	m := md5.New()
	if _, err := io.Copy(m, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(m.Sum(nil)), nil
}

func hashReaderList(r io.Reader, fragmentSizes []int64) (string, []string, error) {
	global := md5.New()
	md5s := make([]string, 0, len(fragmentSizes))
	for _, size := range fragmentSizes {
		local := md5.New()
		if _, err := io.Copy(io.MultiWriter(global, local), io.LimitReader(r, size)); err != nil {
			return "", nil, err
		}
		md5s = append(md5s, hex.EncodeToString(local.Sum(nil)))
	}
	return hex.EncodeToString(global.Sum(nil)), md5s, nil
}

func hashStr(v ...string) string {
	m := md5.New()
	m.Write([]byte(strings.Join(v, "^v^")))
	return hex.EncodeToString(m.Sum(nil))
}

func httpPut(ctx context.Context, url string, reader io.Reader, length int64) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reader)
	if err != nil {
		return err
	}
	request.ContentLength = length
	response, err := new(http.Client).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}
	return errors.New(string(data))
}

func getFragmentSize(totalSize int64, fragmentSize int64) []int64 {
	num := totalSize / fragmentSize
	sizes := make([]int64, num, num+1)
	for i := 0; i < len(sizes); i++ {
		sizes[i] = fragmentSize
	}
	if totalSize%fragmentSize != 0 {
		sizes = append(sizes, totalSize-num*fragmentSize)
	}
	return sizes
}
