package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"path"
)

type OSS struct {
	p *ws.PostApi
}

func NewOSS(p *ws.PostApi) *OSS {
	return &OSS{p: p}
}

func (o *OSS) tencentOssCredentials() (*server_api_params.TencentCloudStorageCredentialRespData, error) {
	req := server_api_params.TencentCloudStorageCredentialReq{OperationID: utils.OperationIDGenerator()}
	commData, err := o.p.PostReturn(constant.TencentCloudStorageCredentialRouter, req)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var resp server_api_params.TencentCloudStorageCredentialResp
	err = mapstructure.Decode(commData.Data, resp.Data)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &resp.Data, nil
}

func (o *OSS) UploadImage(filePath string, callback SendMsgCallBack, operationID string) (string, string) {
	return o.uploadObj(filePath, "img", callback, operationID)

}

func (o *OSS) UploadSound(filePath string, callback SendMsgCallBack, operationID string) (string, string) {
	return o.uploadObj(filePath, "", callback, operationID)
}

func (o *OSS) UploadVideo(videoPath, snapshotPath string, callback SendMsgCallBack, operationID string) (string, string, string, string) {
	videoURL, videoUUID := o.uploadObj(videoPath, "", callback, operationID)
	snapshotURL, snapshotUUID := o.uploadObj(snapshotPath, "img", callback, operationID)
	return snapshotURL, snapshotUUID, videoURL, videoUUID
}

func (o *OSS) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
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

func (o *OSS) uploadObj(filePath string, fileType string, callback SendMsgCallBack, operationID string) (string, string) {
	ossResp, err := o.tencentOssCredentials()
	common.CheckAnyErr(callback, 1000, err, operationID)
	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Credentials.TmpSecretID,
			SecretKey:    ossResp.Credentials.TmpSecretKey,
			SessionToken: ossResp.Credentials.SessionToken,
		},
	})
	if client == nil {
		err := errors.New("client == nil")
		common.CheckAnyErr(callback, 1000, err, operationID)
	}

	newName, contentType, err := o.getNewFileNameAndContentType(filePath, fileType)
	common.CheckAnyErr(callback, 10001, err, operationID)
	var lis = &selfListener{}
	lis.SendMsgCallBack = callback
	opt := &cos.ObjectPutOptions{ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{Listener: lis}}
	if fileType == "img" {
		opt.ContentType = contentType
	}
	_, err = client.Object.PutFromFile(context.Background(), newName, filePath, opt)
	common.CheckAnyErr(callback, 10002, err, operationID)
	return dir + "/" + newName, newName
}

type selfListener struct {
	SendMsgCallBack
}

func (l *selfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressDataEvent:
		if event.ConsumedBytes == event.TotalBytes {
			l.SendMsgCallBack.OnProgress(int((event.ConsumedBytes - 1) * 100 / event.TotalBytes))
		} else {
			l.SendMsgCallBack.OnProgress(int(event.ConsumedBytes * 100 / event.TotalBytes))
		}
	case cos.ProgressFailedEvent:
	}
}
