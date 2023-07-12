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

func NewFile(database db_interface.DataBase, loginUserID string) *File {
	return &File{database: database, loginUserID: loginUserID, lock: &sync.Mutex{}, uploading: make(map[string]func())}
}

type File struct {
	database    db_interface.DataBase
	loginUserID string
	lock        sync.Locker
	partLimit   *third.PartLimitResp
	uploading   map[string]func()
}

func (f *File) cleanPartLimit() {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.partLimit = nil
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

func (f *File) partSize(ctx context.Context, size int64) (int64, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.partLimit == nil {
		resp, err := util.CallApi[third.PartLimitResp](ctx, constant.ObjectPartLimit, &third.PartLimitReq{})
		if err != nil {
			return 0, err
		}
		f.partLimit = resp
	}
	if size <= 0 {
		return 0, errors.New("size must be greater than 0")
	}
	if size > f.partLimit.MaxPartSize*int64(f.partLimit.MaxNumSize) {
		return 0, fmt.Errorf("size must be less than %db", f.partLimit.MaxPartSize*int64(f.partLimit.MaxNumSize))
	}
	if size <= f.partLimit.MinPartSize*int64(f.partLimit.MaxNumSize) {
		return f.partLimit.MinPartSize, nil
	}
	partSize := size / int64(f.partLimit.MaxNumSize)
	if size%int64(f.partLimit.MaxNumSize) != 0 {
		partSize++
	}
	return partSize, nil
}

func (f *File) accessURL(ctx context.Context, req *third.AccessURLReq) (*third.AccessURLResp, error) {
	return util.CallApi[third.AccessURLResp](ctx, constant.ObjectAccessURL, req)
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

func (f *File) partMD5(parts []string) string {
	s := strings.Join(parts, ",")
	md5Sum := md5.Sum([]byte(s))
	return hex.EncodeToString(md5Sum[:])
}

func (f *File) getUploadInfo(ctx context.Context, fileMd5 string, partSize int64) error {

	return nil
}

func (f *File) getUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	return f.initiateMultipartUploadResp(ctx, req)
	//partNum := req.Size / req.PartSize
	//if req.Size%req.PartSize != 0 {
	//	partNum++
	//}
	//var uploadPartIndexs []int32
	//dbUpload, err := f.database.GetUpload(ctx, req.Hash)
	//if err == nil {
	//	uploadPartIndexs, err = f.database.GetUploadPart(ctx, req.Hash)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if partNum <= 1 || len(uploadPartIndexs) == 0 || dbUpload.ExpireTime-3600*1000 < time.Now().UnixMilli() {
	//		if err := f.database.DeleteUpload(ctx, req.Hash); err != nil {
	//			return nil, err
	//		}
	//		dbUpload = nil
	//		uploadPartIndexs = nil
	//	}
	//} else {
	//	dbUpload = nil
	//	log.ZError(ctx, "get upload db", err, "pratsMd5", req.Hash)
	//}
	//if dbUpload == nil {
	//	resp, err := f.initiateMultipartUploadResp(ctx, req)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if resp.Upload != nil {
	//		if err := f.database.InsertUpload(ctx, &model_struct.Upload{
	//			PartHash:   req.Hash,
	//			UploadID:   resp.Upload.UploadID,
	//			ExpireTime: resp.Upload.ExpireTime,
	//			CreateTime: time.Now().UnixMilli(),
	//		}); err != nil {
	//			return nil, err
	//		}
	//	}
	//	return resp, nil
	//}
	//partIndexs := make([]bool, partNum)
	//for _, partIndex := range uploadPartIndexs {
	//	partIndexs[partIndex-1] = true
	//}
	//partNumbers := make([]int32, 0, partNum)
	//for partIndex, ok := range partIndexs {
	//	if !ok {
	//		partNumbers = append(partNumbers, int32(partIndex+1))
	//	}
	//}
	//resp := &third.InitiateMultipartUploadResp{
	//	Upload: &third.UploadInfo{
	//		UploadID:   dbUpload.UploadID,
	//		PartSize:   req.PartSize,
	//		ExpireTime: dbUpload.ExpireTime,
	//		Sign:       &third.AuthSignParts{},
	//	},
	//}
	//if len(partNumbers) > 0 {
	//	authSignResp, err := f.authSign(ctx, &third.AuthSignReq{
	//		UploadID:    dbUpload.UploadID,
	//		PartNumbers: partNumbers,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	resp.Upload.Sign.Url = authSignResp.Url
	//	resp.Upload.Sign.Query = authSignResp.Query
	//	resp.Upload.Sign.Header = authSignResp.Header
	//	resp.Upload.Sign.Parts = authSignResp.Parts
	//}
	//return resp, nil
}

//func (f *File) doHttpPut(ctx context.Context, seeker io.ReadSeeker, upload *third.UploadInfo, index int, partSize int64, size int64) error {
//	sign := upload.Sign
//	item := upload.Sign.Parts[index]
//	num := size / partSize
//	if size%partSize != 0 {
//		num++
//	}
//	currentSize := partSize
//
//
//
//
//
//
//	md5Reader := NewMd5Reader(io.LimitReader(seeker, partSize))
//
//	rawURL := sign.Parts[i].Url
//	if rawURL == "" {
//		rawURL = sign.Url
//	}
//	if len(sign.Query)+len(item.Query) > 0 {
//		u, err := url.Parse(rawURL)
//		if err != nil {
//			return nil, err
//		}
//		query := u.Query()
//		for i := range sign.Query {
//			v := sign.Query[i]
//			query[v.Key] = v.Values
//		}
//		for i := range item.Query {
//			v := item.Query[i]
//			query[v.Key] = v.Values
//		}
//		u.RawQuery = query.Encode()
//		rawURL = u.String()
//	}
//	req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, md5Reader)
//	if err != nil {
//		return nil, err
//	}
//	for i := range sign.Header {
//		v := sign.Header[i]
//		req.Header[v.Key] = v.Values
//	}
//	for i := range item.Header {
//		v := item.Header[i]
//		req.Header[v.Key] = v.Values
//	}
//	if partNum == i+1 {
//		req.ContentLength = size % partSize
//	} else {
//		req.ContentLength = partSize
//	}
//	body, resp, err := f.doHttpReq(req)
//	if err != nil {
//		return nil, err
//	}
//	if resp.StatusCode/200 != 1 {
//		return nil, fmt.Errorf("upload part %d failed, status code %d, body %s", i, resp.StatusCode, string(body))
//	}
//	if md5v := md5Reader.Md5(); md5v != partMd5s[i] {
//		return nil, fmt.Errorf("upload part %d failed, md5 not match, expect %s, got %s", i, partMd5s[i], md5v)
//	}
//	cb.UploadPartComplete(int(item.PartNumber), partSize, partMd5s[i])
//	return nil
//}

func (f *File) UploadFile(ctx context.Context, req *UploadFileReq, cb UploadFileCallback) (*UploadFileResp, error) {
	if cb == nil {
		cb = emptyUploadCallback{}
	}
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
	cb.Open(size)
	partSize, err := f.partSize(ctx, size)
	if err != nil {
		return nil, err
	}
	partNum := int(size / partSize)
	if size%partSize != 0 {
		partNum++
	}
	cb.PartSize(partSize, partNum)
	partMd5s := make([]string, partNum)
	buf := make([]byte, 1024)
	fileMd5 := md5.New()
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
				fileMd5.Write(buf[:n])
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		partMd5s[i] = hex.EncodeToString(h.Sum(nil))
		if partNum == i+1 {
			cb.HashPartProgress(i, size%partSize, partMd5s[i])
		} else {
			cb.HashPartProgress(i, partSize, partMd5s[i])
		}
	}
	partMd5Val := f.partMD5(partMd5s)
	fileMd5Val := hex.EncodeToString(fileMd5.Sum(nil))
	cb.HashPartComplete(f.partMD5(partMd5s), fileMd5Val)
	if _, err := file.Seek(io.SeekStart, 0); err != nil {
		return nil, err
	}
	upload, err := f.getUpload(ctx, &third.InitiateMultipartUploadReq{
		Hash:        partMd5Val,
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
		cb.Complete(size, upload.Url, 0)
		return &UploadFileResp{
			URL: upload.Url,
		}, nil
	}
	if upload.Upload.PartSize != partSize {
		f.cleanPartLimit()
		return nil, fmt.Errorf("part size not match, expect %d, got %d", partSize, upload.Upload.PartSize)
	}
	var uploadSize int64
	for i := 0; i < len(upload.Upload.Sign.Parts); i++ {
		sign := upload.Upload.Sign
		item := upload.Upload.Sign.Parts[i]
		var currentPartSize int64
		if partNum == i+1 {
			currentPartSize = size % partSize
		} else {
			currentPartSize = partSize
		}
		md5Reader := NewMd5Reader(io.LimitReader(file, partSize))
		rawURL := sign.Parts[i].Url
		if rawURL == "" {
			rawURL = sign.Url
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
		req.ContentLength = currentPartSize
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
		uploadSize += currentPartSize
		cb.UploadComplete(size, uploadSize)
		cb.UploadPartComplete(int(item.PartNumber), partSize, partMd5s[i])
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
	cb.Complete(size, resp.Url, 1)
	return &UploadFileResp{
		URL: resp.Url,
	}, nil
}
