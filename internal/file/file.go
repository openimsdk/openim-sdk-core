// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"io"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"os"
	"sync"
)

type PutArgs struct {
	PutID       string
	Filepath    string
	Name        string
	Hash        string
	ContentType string
	ValidTime   int64
}

type PutResp struct {
	URL string
}

func NewFile(dataBase db_interface.DataBase, loginUserID string) *File {
	return &File{loginUserID: loginUserID, lock: new(sync.Mutex), updating: make(map[string]func())}
}

type File struct {
	loginUserID string
	lock        sync.Locker
	updating    map[string]func()
}

func (f *File) apiApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error) {
	return util.CallApi[third.ApplyPutResp](ctx, constant.FileApplyPutRouter, req)
}

func (f *File) apiConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (*third.ConfirmPutResp, error) {
	return util.CallApi[third.ConfirmPutResp](ctx, constant.FileConfirmPutRouter, req)
}

func (f *File) apiGetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error) {
	return util.CallApi[third.GetPutResp](ctx, constant.FileGetPutRouter, req)
}

func (f *File) rePutFilePath(ctx context.Context, req *PutArgs, cb PutFileCallback) (*PutResp, error) {
	file, err := os.Open(req.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	cb.Open(info.Size())
	return f.rePutFile(ctx, file, info.Size(), req, cb)
}

func (f *File) rePutFile(ctx context.Context, file *os.File, size int64, req *PutArgs, cb PutFileCallback) (*PutResp, error) {
	if req.Hash == "" {
		var err error
		req.Hash, err = hashReader(NewReader(ctx, file, size, cb.HashProgress))
		if err != nil {
			return nil, err
		}
		cb.HashComplete(req.Hash, size)
		if _, err := file.Seek(io.SeekStart, 0); err != nil {
			return nil, err
		}
	} else {
		if v, err := hex.DecodeString(req.Hash); err != nil {
			return nil, err
		} else if len(v) != md5.Size {
			return nil, fmt.Errorf("hash length error")
		}
	}
	applyPutResp, err := f.apiApplyPut(ctx, &third.ApplyPutReq{PutID: req.PutID, Name: req.Name, ContentType: req.ContentType, ValidTime: req.ValidTime, Hash: req.Hash, Size: size})
	if err != nil {
		return nil, err
	}
	if applyPutResp.Url != "" {
		cb.PutStart(size, size)
		cb.PutComplete(size, 0)
		return &PutResp{URL: applyPutResp.Url}, nil
	}
	req.PutID = applyPutResp.PutID
	cb.PutStart(0, size)
	fragments := getFragmentSize(size, applyPutResp.FragmentSize)
	if len(fragments) != len(applyPutResp.PutURLs) {
		return nil, fmt.Errorf("get fragment size error local %d server %d", len(fragments), len(applyPutResp.PutURLs))
	}
	var initSize int64
	for i, url := range applyPutResp.PutURLs {
		put := NewReader(ctx, io.LimitReader(file, fragments[i]), size, func(current, total int64) {
			cb.PutProgress(initSize, current+initSize, size)
		})
		if err := httpPut(ctx, url, put, fragments[i]); err != nil {
			return nil, err
		}
		initSize += fragments[i]
	}
	cb.PutProgress(size, size, size) // 100%
	confirmPutResp, err := f.apiConfirmPut(ctx, &third.ConfirmPutReq{PutID: applyPutResp.PutID})
	if err != nil {
		return nil, err
	}
	cb.PutComplete(size, 1)
	return &PutResp{URL: confirmPutResp.Url}, nil
}

func (f *File) putFile(ctx context.Context, req *PutArgs, cb PutFileCallback) (*PutResp, error) {
	if req.PutID == "" {
		return f.rePutFilePath(ctx, req, cb) // 没有putID
	}
	resp, err := f.apiGetPut(ctx, &third.GetPutReq{PutID: req.PutID})
	if errs.ErrRecordNotFound.Is(err) {
		return f.rePutFilePath(ctx, req, cb) // 服务端不存在，重新上传
	} else if errs.ErrFileUploadedExpired.Is(err) {
		return f.rePutFilePath(ctx, req, cb) // 上传时间过期
	} else if err != nil {
		return nil, err // 其他错误
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
	fragmentSizes := getFragmentSize(resp.Size, resp.FragmentSize)
	hash, md5s, err := hashReaderList(NewReader(ctx, file, info.Size(), cb.HashProgress), fragmentSizes)
	if err != nil {
		return nil, err
	}
	if resp.Size != info.Size() || resp.Hash != hash {
		return nil, errors.New("file size or hash error")
	}
	if len(md5s) != len(resp.Fragments) {
		return nil, fmt.Errorf("get fragment size error local %d server %d", len(md5s), len(resp.Fragments))
	}
	var putSize int64               // 已上传的大小
	puts := make([]bool, len(md5s)) // 已上传的片段
	for i, fragment := range resp.Fragments {
		if fragment.Hash == md5s[i] {
			puts[i] = true
			putSize += fragmentSizes[i]
		}
	}
	var readSize int64 // 已读取的大小
	for i, fragment := range resp.Fragments {
		if puts[i] {
			readSize += fragmentSizes[i]
			continue
		}
		if _, err := file.Seek(io.SeekStart, int(readSize)); err != nil {
			return nil, err
		}
		reader := NewReader(ctx, io.LimitReader(file, fragmentSizes[i]), info.Size(), func(current, total int64) {
			cb.PutProgress(putSize, current+putSize, info.Size())
		})
		if err := httpPut(ctx, fragment.Url, reader, fragmentSizes[i]); err != nil {
			return nil, err
		}
		putSize += fragmentSizes[i]
		readSize += fragmentSizes[i]
	}
	cb.PutProgress(info.Size(), info.Size(), info.Size())
	confirmPutResp, err := f.apiConfirmPut(ctx, &third.ConfirmPutReq{PutID: req.PutID})
	if err != nil {
		return nil, err
	}
	cb.PutComplete(info.Size(), 2)
	return &PutResp{URL: confirmPutResp.Url}, nil
}

func (f *File) PutFile(ctx context.Context, req *PutArgs, cb PutFileCallback) (*PutResp, error) {
	if req.PutID == "" {
		return nil, fmt.Errorf("put id is empty")
	}
	f.lock.Lock()
	if _, ok := f.updating[req.PutID]; ok {
		f.lock.Unlock()
		return nil, fmt.Errorf("put id is uploading")
	}
	done := ctx.Done()
	ctx, cancel := context.WithCancel(ctx)
	if done != nil {
		go func() {
			<-done
			cancel()
		}()
	}
	f.updating[req.PutID] = cancel
	f.lock.Unlock()
	defer func(putID string) {
		f.lock.Lock()
		delete(f.updating, putID)
		f.lock.Unlock()
	}(req.PutID)
	return f.putFile(ctx, req, cb)
}

func (f *File) Cancel(ctx context.Context, putID string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	cancel, ok := f.updating[putID]
	if ok {
		delete(f.updating, putID)
	}
	cancel()
}
