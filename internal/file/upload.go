package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"io"
	"net/http"
	"net/url"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"os"
	"strings"
	"sync"
)

type UploadFileReq struct {
	Filepath    string `json:"filepath"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Cause       string `json:"cause"`
}

type UploadFileResp struct {
	URL string `json:"url"`
}

func NewFile(dataBase db_interface.DataBase, loginUserID string) *File {
	return &File{loginUserID: loginUserID, lock: new(sync.Mutex), uploading: make(map[string]func())}
}

type File struct {
	loginUserID string
	lock        sync.Locker
	uploading   map[string]func()
}

func (f *File) initiateMultipartUploadResp(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	return util.CallApi[third.InitiateMultipartUploadResp](ctx, constant.ObjectInitiateMultipartUpload, req)
}

func (f *File) authSign(ctx context.Context, req *third.AuthSignReq) (*third.AuthSignResp, error) {
	return util.CallApi[third.AuthSignResp](ctx, constant.ObjectAuthSign, req)
}

func (f *File) completeMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	return util.CallApi[third.CompleteMultipartUploadResp](ctx, constant.ObjectCompleteMultipartUpload, req)
}

func (f *File) partSize(ctx context.Context, req *third.PartSizeReq) (*third.PartSizeResp, error) {
	return util.CallApi[third.PartSizeResp](ctx, constant.ObjectPartSize, req)
}

func (f *File) getPartSize(ctx context.Context, size int64) (int64, error) {
	resp, err := f.partSize(ctx, &third.PartSizeReq{Size: size})
	if err != nil {
		return 0, err
	}
	return resp.Size, nil
}

func (f *File) doHttpReq(req *http.Request) ([]byte, *http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return data, resp, nil
}

func (f *File) UploadFile(ctx context.Context, req *UploadFileReq) (*UploadFileResp, error) {
	if req.Name == "" {
		return nil, errors.New("name is empty")
	}
	if req.Name[0] == '/' {
		req.Name = req.Name[1:]
	}
	if prefix := f.loginUserID + "/"; !strings.HasPrefix(req.Name, prefix) {
		req.Name = prefix + req.Name
	}
	file, err := os.Open(req.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	partSize, err := f.getPartSize(ctx, size)
	if err != nil {
		return nil, err
	}
	partNum := int(size / partSize)
	if size%partSize != 0 {
		partNum++
	}
	partMd5s := make([]string, partNum)
	buf := make([]byte, 1024)
	for i := 0; i < partNum; i++ {
		h := md5.New()
		r := io.LimitReader(file, partSize)
		for {
			n, err := r.Read(buf)
			if err == nil {
				if req.ContentType == "" {
					req.ContentType = http.DetectContentType(buf[:n])
				}
				h.Write(buf[:n])
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		partMd5s[i] = hex.EncodeToString(h.Sum(nil))
	}
	if _, err := file.Seek(io.SeekStart, 0); err != nil {
		return nil, err
	}
	upload, err := f.initiateMultipartUploadResp(ctx, &third.InitiateMultipartUploadReq{
		Hash:        Md5Str(strings.Join(partMd5s, ",")),
		Size:        size,
		PartSize:    partSize,
		MaxParts:    -1,
		Cause:       req.Cause,
		Name:        req.Name,
		ContentType: req.ContentType,
	})
	if err != nil {
		return nil, err
	}
	if upload.Upload == nil {
		return &UploadFileResp{
			URL: upload.Url,
		}, nil
	}
	if upload.Upload.PartSize != partSize {
		return nil, fmt.Errorf("part size not match, expect %d, got %d", partSize, upload.Upload.PartSize)
	}
	for i := 0; i < partNum; i++ {
		sign := upload.Upload.Sign
		item := upload.Upload.Sign.Parts[i]
		md5Reader := NewMd5Reader(io.LimitReader(file, partSize))
		rawURL := sign.Parts[i].Url
		if rawURL == "" {
			rawURL = upload.Url
		}
		if len(sign.Query)+len(item.Query) > 0 {
			u, err := url.Parse(rawURL)
			if err != nil {
				return nil, err
			}
			query := u.Query()
			for i := range sign.Query {
				v := sign.Query[i]
				query[v.Key] = v.Values
			}
			for i := range item.Query {
				v := item.Query[i]
				query[v.Key] = v.Values
			}
			u.RawQuery = query.Encode()
			rawURL = u.String()
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, md5Reader)
		if err != nil {
			return nil, err
		}
		for i := range sign.Header {
			v := sign.Header[i]
			req.Header[v.Key] = v.Values
		}
		for i := range item.Header {
			v := item.Header[i]
			req.Header[v.Key] = v.Values
		}
		if partNum == i+1 {
			req.ContentLength = size % partSize
		} else {
			req.ContentLength = partSize
		}
		body, resp, err := f.doHttpReq(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode/200 != 1 {
			return nil, fmt.Errorf("upload part %d failed, status code %d, body %s", i, resp.StatusCode, string(body))
		}
		if md5v := md5Reader.Md5(); md5v != partMd5s[i] {
			return nil, fmt.Errorf("upload part %d failed, md5 not match, expect %s, got %s", i, partMd5s[i], md5v)
		}
	}
	resp, err := f.completeMultipartUpload(ctx, &third.CompleteMultipartUploadReq{
		UploadID:    upload.Upload.UploadID,
		Parts:       partMd5s,
		Name:        req.Name,
		ContentType: req.ContentType,
		Cause:       req.Cause,
	})
	if err != nil {
		return nil, err
	}
	return &UploadFileResp{
		URL: resp.Url,
	}, nil
}
