package obj_storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"

	"math/rand"
	ws "open_im_sdk/internal/interaction"
	//	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"path"
	"time"
)

// OSS 阿里云对象存储
type OSS struct {
	p *ws.PostApi
}

func (c *OSS) UploadImageByBuffer(buffer *bytes.Buffer, size int64, imageType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *OSS) UploadSoundByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *OSS) UploadFileByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (c *OSS) UploadVideoByBuffer(videoBuffer, snapshotBuffer *bytes.Buffer, videoSize, snapshotSize int64, videoType string, onProgressFun func(int)) (string, string, string, string, error) {
	panic("implement me")
}

func NewOSS(p *ws.PostApi) *OSS {
	return &OSS{p: p}
}

func (c *OSS) aliOSSCredentials(filename string, fileType string) (*server_api_params.OSSCredentialRespData, error) {
	req := server_api_params.OSSCredentialReq{
		OperationID: utils.OperationIDGenerator(),
		Filename:    filename,
		FileType:    fileType,
	}
	var resp server_api_params.OSSCredentialResp
	err := c.p.PostReturn(constant.AliOSSCredentialRouter, req, &resp.OssData)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &resp.OssData, nil
}

func (c *OSS) UploadImage(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "img", onProgressFun)

}

func (c *OSS) UploadSound(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "", onProgressFun)
}

func (c *OSS) UploadFile(filePath string, onProgressFun func(int)) (string, string, error) {
	return c.uploadObj(filePath, "", onProgressFun)
}

func (c *OSS) UploadVideo(videoPath, snapshotPath string, onProgressFun func(int)) (string, string, string, string, error) {
	videoURL, videoUUID, err := c.uploadObj(videoPath, "", onProgressFun)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	snapshotURL, snapshotUUID, err := c.uploadObj(snapshotPath, "img", onProgressFun)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	return snapshotURL, snapshotUUID, videoURL, videoUUID, nil
}

func (c *OSS) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
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

// uploadObj 上传对象 onProgressFun 返回值代表百分比 -1表示失败
func (c *OSS) uploadObj(filePath string, fileType string, onProgressFun func(int)) (string, string, error) {
	OSSResp, err := c.aliOSSCredentials(filePath, fileType)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	log.Info("upload Endpoint: ", OSSResp.Endpoint, "bucket: ", OSSResp.Bucket)
	client, err := oss.New(OSSResp.Endpoint, OSSResp.AccessKeyId, OSSResp.AccessKeySecret, oss.SecurityToken(OSSResp.Token))
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	// 获取存储空间。
	bucket, err := client.Bucket(OSSResp.Bucket)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	newName, _, err := c.getNewFileNameAndContentType(filePath, fileType)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	// 带可选参数的签名直传。请确保设置的ContentType值与在前端使用时设置的ContentType值一致。
	options := []oss.Option{
		oss.Progress(&OssProgressListener{onProgressFun: onProgressFun}), // 进度条
	}
	// 签名直传。
	signedURL, err := bucket.SignURL(newName, oss.HTTPPut, 60)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	err = bucket.PutObjectFromFileWithURL(signedURL, filePath, options...)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	return OSSResp.FinalHost + "/" + newName, newName, nil
}

// OssProgressListener 定义进度条监听器。
type OssProgressListener struct {
	onProgressFun func(int)
}

// ProgressChanged 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferDataEvent:
		listener.onProgressFun(int((event.ConsumedBytes - 1) * 100 / event.TotalBytes))
	case oss.TransferCompletedEvent:
		listener.onProgressFun(100)
	case oss.TransferFailedEvent:
		listener.onProgressFun(-1)
	default:
	}
}
