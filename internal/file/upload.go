// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/log"
)

type UploadFileReq struct {
	Filepath    string `json:"filepath"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Cause       string `json:"cause"`
	Uuid        string `json:"uuid"`
}

type UploadFileResp struct {
	URL string `json:"url"`
}

type partInfo struct {
	ContentType string
	PartSize    int64
	PartNum     int
	FileMd5     string
	PartMd5     string
	PartSizes   []int64
	PartMd5s    []string
}

func NewFile(database db_interface.DataBase, loginUserID string) *File {
	return &File{database: database, loginUserID: loginUserID, confLock: &sync.Mutex{}, mapLocker: &sync.Mutex{}, uploading: make(map[string]*lockInfo)}
}

type File struct {
	database    db_interface.DataBase
	loginUserID string
	confLock    sync.Locker
	partLimit   *third.PartLimitResp
	mapLocker   sync.Locker
	uploading   map[string]*lockInfo
}

type lockInfo struct {
	count  int32
	locker sync.Locker
}

func (f *File) lockHash(hash string) {
	f.mapLocker.Lock()
	locker, ok := f.uploading[hash]
	if !ok {
		locker = &lockInfo{count: 0, locker: &sync.Mutex{}}
		f.uploading[hash] = locker
	}
	atomic.AddInt32(&locker.count, 1)
	f.mapLocker.Unlock()
	locker.locker.Lock()
}

func (f *File) unlockHash(hash string) {
	f.mapLocker.Lock()
	locker, ok := f.uploading[hash]
	if !ok {
		f.mapLocker.Unlock()
		return
	}
	if atomic.AddInt32(&locker.count, -1) == 0 {
		delete(f.uploading, hash)
	}
	f.mapLocker.Unlock()
	locker.locker.Unlock()
}

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
	file, err := Open(req)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileSize := file.Size()
	cb.Open(fileSize)
	info, err := f.getPartInfo(ctx, file, fileSize, cb)
	if err != nil {
		return nil, err
	}
	if req.ContentType == "" {
		req.ContentType = info.ContentType
	}
	partSize := info.PartSize
	partSizes := info.PartSizes
	partMd5s := info.PartMd5s
	partMd5Val := info.PartMd5
	if err := file.StartSeek(0); err != nil {
		return nil, err
	}
	f.lockHash(partMd5Val)
	defer f.unlockHash(partMd5Val)
	// 获取上传文件
	bitmap, dbUpload, upload, err := f.getUpload(ctx, &third.InitiateMultipartUploadReq{
		Hash:        partMd5Val,
		Size:        fileSize,
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
		cb.Complete(fileSize, upload.Url, 0)
		return &UploadFileResp{
			URL: upload.Url,
		}, nil
	}
	if upload.Upload.PartSize != partSize {
		f.cleanPartLimit()
		return nil, fmt.Errorf("part fileSize not match, expect %d, got %d", partSize, upload.Upload.PartSize)
	}
	cb.UploadID(upload.Upload.UploadID)
	uploadedSize := fileSize
	uploadParts := make([]*third.SignPart, info.PartNum)
	for _, part := range upload.Upload.Sign.Parts {
		uploadParts[part.PartNumber-1] = part
		uploadedSize -= partSizes[part.PartNumber-1]
	}
	continueUpload := uploadedSize > 0
	for i, currentPartSize := range partSizes {
		md5Reader := NewMd5Reader(io.LimitReader(file, currentPartSize))
		part := uploadParts[i]
		if part == nil {
			if _, err := io.Copy(io.Discard, md5Reader); err != nil {
				return nil, err
			}
		} else {
			reader := NewProgressReader(md5Reader, func(current int64) {
				cb.UploadComplete(fileSize, uploadedSize+current, uploadedSize)
			})
			if err := f.doPut(ctx, http.DefaultClient, upload.Upload.Sign, part, reader, currentPartSize); err != nil {
				return nil, err
			}
			uploadedSize += currentPartSize
		}
		if md5val := md5Reader.Md5(); md5val != partMd5s[i] {
			return nil, fmt.Errorf("upload part %d failed, md5 not match, expect %s, got %s", i, partMd5s[i], md5val)
		}
		if part != nil && dbUpload != nil && bitmap != nil {
			bitmap.Set(int(part.PartNumber - 1))
			dbUpload.UploadInfo = base64.StdEncoding.EncodeToString(bitmap.Serialize())
			if err := f.database.UpdateUpload(ctx, dbUpload); err != nil {
				log.ZError(ctx, "SetUploadPartPush", err, "partMd5Val", partMd5Val, "name", req.Name, "partNumber", part.PartNumber)
			}
		}
		cb.UploadPartComplete(i, currentPartSize, partMd5s[i])
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
	typ := 1
	if continueUpload {
		typ++
	}
	cb.Complete(fileSize, resp.Url, typ)
	if err := f.database.DeleteUpload(ctx, info.PartMd5); err != nil {
		log.ZError(ctx, "DeleteUpload", err, "partMd5Val", info.PartMd5, "name", req.Name)
	}
	return &UploadFileResp{
		URL: resp.Url,
	}, nil
}

func (f *File) cleanPartLimit() {
	f.confLock.Lock()
	defer f.confLock.Unlock()
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

func (f *File) getPartNum(fileSize int64, partSize int64) int {
	partNum := fileSize / partSize
	if fileSize%partSize != 0 {
		partNum++
	}
	return int(partNum)
}

func (f *File) partSize(ctx context.Context, size int64) (int64, error) {
	f.confLock.Lock()
	defer f.confLock.Unlock()
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

func (f *File) getUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*Bitmap, *model_struct.LocalUpload, *third.InitiateMultipartUploadResp, error) {
	partNum := f.getPartNum(req.Size, req.PartSize)
	dbUpload, err := f.database.GetUpload(ctx, req.Hash)
	var bitmap *Bitmap
	if err == nil {
		bitmapBytes, err := base64.StdEncoding.DecodeString(dbUpload.UploadInfo)
		if len(bitmapBytes) == 0 || err != nil || partNum <= 1 || dbUpload.ExpireTime-3600*1000 < time.Now().UnixMilli() {
			if err := f.database.DeleteUpload(ctx, req.Hash); err != nil {
				return nil, nil, nil, err
			}
			dbUpload = nil
		}
		if dbUpload != nil {
			bitmap = ParseBitmap(bitmapBytes, partNum)
		}
	} else {
		dbUpload = nil
		log.ZError(ctx, "get upload db", err, "pratsMd5", req.Hash)
	}
	if dbUpload == nil {
		resp, err := f.initiateMultipartUploadResp(ctx, req)
		if err != nil {
			return nil, nil, nil, err
		}
		if resp.Upload == nil {
			return nil, nil, resp, nil
		}
		bitmap = NewBitmap(partNum)
		if resp.Upload != nil {
			dbUpload = &model_struct.LocalUpload{
				PartHash:   req.Hash,
				UploadID:   resp.Upload.UploadID,
				UploadInfo: base64.StdEncoding.EncodeToString(bitmap.Serialize()),
				ExpireTime: resp.Upload.ExpireTime,
				CreateTime: time.Now().UnixMilli(),
			}
			if err := f.database.InsertUpload(ctx, dbUpload); err != nil {
				return nil, nil, nil, err
			}
		}
		return bitmap, dbUpload, resp, nil
	}
	partNumbers := make([]int32, 0, partNum)
	for i := 0; i < partNum; i++ {
		if !bitmap.Get(i) {
			partNumbers = append(partNumbers, int32(i+1))
		}
	}
	resp := &third.InitiateMultipartUploadResp{
		Upload: &third.UploadInfo{
			UploadID:   dbUpload.UploadID,
			PartSize:   req.PartSize,
			ExpireTime: dbUpload.ExpireTime,
			Sign:       &third.AuthSignParts{},
		},
	}
	if len(partNumbers) > 0 {
		authSignResp, err := f.authSign(ctx, &third.AuthSignReq{
			UploadID:    dbUpload.UploadID,
			PartNumbers: partNumbers,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		resp.Upload.Sign.Url = authSignResp.Url
		resp.Upload.Sign.Query = authSignResp.Query
		resp.Upload.Sign.Header = authSignResp.Header
		resp.Upload.Sign.Parts = authSignResp.Parts
	}
	return bitmap, dbUpload, resp, nil
}

func (f *File) doPut(ctx context.Context, client *http.Client, sign *third.AuthSignParts, part *third.SignPart, reader io.Reader, size int64) error {
	rawURL := part.Url
	if rawURL == "" {
		rawURL = sign.Url
	}
	if len(sign.Query)+len(part.Query) > 0 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		query := u.Query()
		for i := range sign.Query {
			v := sign.Query[i]
			query[v.Key] = v.Values
		}
		for i := range part.Query {
			v := part.Query[i]
			query[v.Key] = v.Values
		}
		u.RawQuery = query.Encode()
		rawURL = u.String()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, reader)
	if err != nil {
		return err
	}
	for i := range sign.Header {
		v := sign.Header[i]
		req.Header[v.Key] = v.Values
	}
	for i := range part.Header {
		v := part.Header[i]
		req.Header[v.Key] = v.Values
	}
	req.ContentLength = size
	log.ZDebug(ctx, "do put req", "url", rawURL, "contentLength", size, "header", req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	log.ZDebug(ctx, "do put resp status", "url", rawURL, "status", resp.Status, "contentLength", resp.ContentLength, "header", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "do put resp body", "url", rawURL, "body", string(body))
	if resp.StatusCode/200 != 1 {
		return fmt.Errorf("PUT %s part %d failed, status code %d, body %s", rawURL, part.PartNumber, resp.StatusCode, string(body))
	}
	return nil
}

func (f *File) getPartInfo(ctx context.Context, r io.Reader, fileSize int64, cb UploadFileCallback) (*partInfo, error) {
	partSize, err := f.partSize(ctx, fileSize)
	if err != nil {
		return nil, err
	}
	partNum := int(fileSize / partSize)
	if fileSize%partSize != 0 {
		partNum++
	}
	cb.PartSize(partSize, partNum)
	partSizes := make([]int64, partNum)
	for i := 0; i < partNum; i++ {
		partSizes[i] = partSize
	}
	partSizes[partNum-1] = fileSize - partSize*(int64(partNum)-1)
	partMd5s := make([]string, partNum)
	buf := make([]byte, 1024*8)
	fileMd5 := md5.New()
	var contentType string
	for i := 0; i < partNum; i++ {
		h := md5.New()
		r := io.LimitReader(r, partSize)
		for {
			if n, err := r.Read(buf); err == nil {
				if contentType == "" {
					contentType = http.DetectContentType(buf[:n])
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
		cb.HashPartProgress(i, partSizes[i], partMd5s[i])
	}
	partMd5Val := f.partMD5(partMd5s)
	fileMd5val := hex.EncodeToString(fileMd5.Sum(nil))
	cb.HashPartComplete(f.partMD5(partMd5s), hex.EncodeToString(fileMd5.Sum(nil)))
	return &partInfo{
		ContentType: contentType,
		PartSize:    partSize,
		PartNum:     partNum,
		FileMd5:     fileMd5val,
		PartMd5:     partMd5Val,
		PartSizes:   partSizes,
		PartMd5s:    partMd5s,
	}, nil
}
