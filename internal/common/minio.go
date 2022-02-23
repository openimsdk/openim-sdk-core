package common

import (
	"context"
	"errors"
	"fmt"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"math/rand"
	"net/url"
	_"net/url"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"path"
	"time"
)



type Minio struct {
	p *ws.PostApi
}

func NewMinio(p *ws.PostApi) *Minio {
	return &Minio{p:p}
}


func (m *Minio) getMinioCredentials() (*server_api_params.MinioStorageCredentialResp, error){
	req := server_api_params.MinioStorageCredentialReq{OperationID: utils.OperationIDGenerator()}
	var resp server_api_params.MinioStorageCredentialResp
	err := m.p.PostReturn(constant.MinioStorageCredentialRouter, req, &resp)
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		return &resp, utils.Wrap(err, "")
	}
	return &resp, nil
}

func (m *Minio) upload(filePath, fileType string, onProgressFun func(int)) (string, string, error){
	minioResp, err := m.getMinioCredentials()
	endPoint, err  := url.Parse(minioResp.StsEndpointURL)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "url parse failed", err.Error())
		return "", "", err
	}
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "new minio client failed", err.Error())
	}
	newName, newType, err := m.getNewFileNameAndContentType(filePath, fileType)
	client, err := minio.New(endPoint.Host,  &minio.Options{
		Creds:        credentials.NewStaticV4(minioResp.AccessKeyID, minioResp.SecretAccessKey, minioResp.SessionToken),
		Secure:       false,
	})
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "generate filename and filetype failed", err.Error())
		return "", "", utils.Wrap(err, "")
	}
	_, err = client.FPutObject(context.Background(), minioResp.BucketName, newName, filePath, minio.PutObjectOptions{ContentType:newType})
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "FPutObject failed", err.Error())
		return "", "", utils.Wrap(err, "")
	}
	// fake callback
	onProgressFun(100)
	presignedURL := endPoint.String() + "/" + minioResp.BucketName + "/" + newName
	return presignedURL, newName, nil
}


func (m *Minio) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
	suffix := path.Ext(filePath)
	if len(suffix) == 0 {
		return "", "", utils.Wrap(errors.New("no suffix "), filePath)
	}
	newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
	contentType := ""
	if fileType == "img" {
		contentType = "image/" + suffix[1:]
	}
	return newName, contentType, nil
}

func (m *Minio) UploadImage(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "img	", onProgressFun)
}


func (m *Minio) UploadSound(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "", onProgressFun)
}

func (m *Minio) UploadFile(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "", onProgressFun)
}

func (m *Minio) UploadVideo(videoPath, snapshotPath string, onProgressFun func(int)) (string, string, string, string, error) {
	videoURL, videoName, err := m.upload(videoPath, "", onProgressFun)
	snapshotURL, snapshotUUID, err :=  m.upload(snapshotPath, "img", onProgressFun)
	if err != nil {
		return snapshotURL, snapshotUUID, videoURL, videoName, utils.Wrap(err, "")
	}
	return snapshotURL, snapshotUUID, videoURL, videoName, nil
}

