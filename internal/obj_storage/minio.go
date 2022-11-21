package obj_storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"math/rand"
	"net/url"
	_ "net/url"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"os"
	"path"
	"time"
)

type Minio struct {
	p *ws.PostApi
}

func (m *Minio) UploadImageByBuffer(buffer *bytes.Buffer, size int64, imageType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *Minio) UploadSoundByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *Minio) UploadFileByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *Minio) UploadVideoByBuffer(videoBuffer, snapshotBuffer *bytes.Buffer, videoSize, snapshotSize int64, videoType string, onProgressFun func(int)) (string, string, string, string, error) {
	panic("implement me")
}

func NewMinio(p *ws.PostApi) *Minio {
	return &Minio{p: p}
}

func (m *Minio) getMinioCredentials() (*server_api_params.MinioStorageCredentialResp, error) {
	req := server_api_params.MinioStorageCredentialReq{OperationID: utils.OperationIDGenerator()}
	var resp server_api_params.MinioStorageCredentialResp
	err := m.p.PostReturn(constant.MinioStorageCredentialRouter, req, &resp)
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error(), resp, req)
		return &resp, utils.Wrap(err, "")
	}
	return &resp, nil
}

func (m *Minio) upload(filePath, fileType string, onProgressFun func(int)) (string, string, error) {
	minioResp, err := m.getMinioCredentials()
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "getMinioCredentials from server failed, please check server log", err.Error(), "resp: ", *minioResp)
		return "", "", utils.Wrap(err, "")
	}
	log.NewInfo("", utils.GetSelfFuncName(), "recv minio credentials", *minioResp)
	endPoint, err := url.Parse(minioResp.StsEndpointURL)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "url parse failed, pleace check config/config.yaml", err.Error())
		return "", "", utils.Wrap(err, "")
	}
	newName, newType, err := m.getNewFileNameAndContentType(filePath, fileType)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "getNewFileNameAndContentType failed", err.Error(), filePath, fileType)
		return "", "", utils.Wrap(err, "")
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(minioResp.AccessKeyID, minioResp.SecretAccessKey, minioResp.SessionToken),
	}
	switch endPoint.Scheme {
	case "http":
		opts.Secure = false
	case "https":
		opts.Secure = true
	}
	client, err := minio.New(endPoint.Host, opts)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "generate filename and filetype failed", err.Error(), endPoint.Host)
		return "", "", utils.Wrap(err, "")
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return "", "", utils.Wrap(err, "os stat failed")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	defer file.Close()

	progressBar := NewUploadProgress(fi.Size(), onProgressFun)
	//_, err = client.FPutObject(context.Background(), minioResp.BucketName, newName, filePath, minio.PutObjectOptions{ContentType: newType,
	//	Progress: progressBar, RetainUntilDate:time.Now().Add(time.Duration(minioResp.StorageTime)*time.Hour*24), Mode: "", Rente})
	_, err = client.PutObject(context.Background(), minioResp.BucketName, newName, file, fi.Size(), minio.PutObjectOptions{ContentType: newType,
		Progress: progressBar})
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "FPutObject failed", err.Error(), newName, filePath, newType)
		return "", "", utils.Wrap(err, "")
	}
	presignedURL := endPoint.String() + "/" + minioResp.BucketName + "/" + newName
	return presignedURL, newName, nil
}

func (m *Minio) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
	suffix := path.Ext(filePath)
	newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
	contentType := ""
	if fileType == "img" {
		if len(suffix) == 0 {
			return "", "", utils.Wrap(errors.New("no suffix "), filePath)
		} else {
			contentType = "image/" + suffix[1:]
		}
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
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	snapshotURL, snapshotUUID, err := m.upload(snapshotPath, "img", nil)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	return snapshotURL, snapshotUUID, videoURL, videoName, nil
}

// NewUploadProgress: Create a *UploadProgress object
func NewUploadProgress(total int64, f func(int)) *UploadProgress {
	return &UploadProgress{total: total, current: 0, percent: 0, callbackFunc: f}
}

type UploadProgress struct {
	total        int64
	current      int64
	percent      int
	callbackFunc func(int)
}

func (progress *UploadProgress) Read(b []byte) (int, error) {
	n := int64(len(b))
	progress.current += n
	percent := int(float64(progress.current) * 100 / float64(progress.total))
	if percent != progress.percent && percent%5 == 0 {
		progress.percent = percent
		if progress.callbackFunc != nil {
			progress.callbackFunc(progress.percent)
		}
	}
	return int(n), nil
}
