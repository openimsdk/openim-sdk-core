package obj_storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"

	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	"net/http"
	"net/url"
	ws "open_im_sdk/internal/interaction"
	//	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"path"
	"time"
)

type COS struct {
	p *ws.PostApi
}

func (c *COS) UploadImageByBuffer(buffer *bytes.Buffer, size int64, imageType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *COS) UploadSoundByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *COS) UploadFileByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *COS) UploadVideoByBuffer(videoBuffer, snapshotBuffer *bytes.Buffer, videoSize, snapshotSize int64, videoType string, onProgressFun func(int)) (string, string, string, string, error) {
	panic("implement me")
}

func NewCOS(p *ws.PostApi) *COS {
	return &COS{p: p}
}

func (c *COS) tencentCOSCredentials() (*server_api_params.TencentCloudStorageCredentialRespData, error) {
	req := server_api_params.TencentCloudStorageCredentialReq{OperationID: utils.OperationIDGenerator()}
	var resp server_api_params.TencentCloudStorageCredentialResp
	err := c.p.PostReturn(constant.TencentCloudStorageCredentialRouter, req, &resp.CosData)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &resp.CosData, nil
}

func (c *COS) UploadImage(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "img", onProgressFun)

}

func (c *COS) UploadSound(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "", onProgressFun)
}

func (c *COS) UploadFile(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "", onProgressFun)
}

func (c *COS) UploadVideo(videoPath, snapshotPath string, onProgressFun func(int)) (string, string, string, string, error) {
	videoURL, videoUUID, err := c.uploadObj(videoPath, "", onProgressFun)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	snapshotURL, snapshotUUID, err := c.uploadObj(snapshotPath, "img", nil)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	return snapshotURL, snapshotUUID, videoURL, videoUUID, nil
}

func (c *COS) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
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

func (c *COS) uploadObj(filePath string, fileType string, onProgressFun func(int)) (string, string, error) {
	COSResp, err := c.tencentCOSCredentials()
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	log.Info("upload ", COSResp.Credentials.SessionToken, "bucket ", COSResp.Credentials.TmpSecretID)
	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", COSResp.Bucket, COSResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     COSResp.Credentials.TmpSecretID,
			SecretKey:    COSResp.Credentials.TmpSecretKey,
			SessionToken: COSResp.Credentials.SessionToken,
		},
	})
	if client == nil {
		err := errors.New("client == nil")
		return "", "", utils.Wrap(err, "")
	}

	newName, contentType, err := c.getNewFileNameAndContentType(filePath, fileType)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	var lis = &selfListener{}
	lis.onProgressFun = onProgressFun
	opt := &cos.ObjectPutOptions{ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{Listener: lis}}
	if fileType == "img" {
		opt.ContentType = contentType
	}
	_, err = client.Object.PutFromFile(context.Background(), newName, filePath, opt)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	return dir + "/" + newName, newName, nil
}

type selfListener struct {
	onProgressFun func(int)
}

func (l *selfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressDataEvent:
		if event.ConsumedBytes == event.TotalBytes {
			if l.onProgressFun != nil {
				l.onProgressFun(int((event.ConsumedBytes - 1) * 100 / event.TotalBytes))
			}

		} else {
			if l.onProgressFun != nil {
				l.onProgressFun(int(event.ConsumedBytes * 100 / event.TotalBytes))
			}
		}
	case cos.ProgressFailedEvent:
	}
}
