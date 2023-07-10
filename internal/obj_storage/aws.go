package obj_storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWS struct {
	p *ws.PostApi
}

func (m *AWS) UploadImageByBuffer(buffer *bytes.Buffer, size int64, imageType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *AWS) UploadSoundByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *AWS) UploadFileByBuffer(buffer *bytes.Buffer, size int64, fileType string, onProgressFun func(int)) (string, string, error) {
	panic("implement me")
}

func (m *AWS) UploadVideoByBuffer(videoBuffer, snapshotBuffer *bytes.Buffer, videoSize, snapshotSize int64, videoType string, onProgressFun func(int)) (string, string, string, string, error) {
	panic("implement me")
}

func NewAWS(p *ws.PostApi) *AWS {
	return &AWS{p: p}
}

func (m *AWS) getAwsCredentials() (server_api_params.AwsStorageCredentialResp, error) {
	req := server_api_params.AwsStorageCredentialReq{OperationID: utils.OperationIDGenerator()}
	var resp server_api_params.AwsStorageCredentialResp
	err := m.p.PostReturn(constant.AwsStorageCredentialRouter, req, &resp)
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error(), resp, req)
		return resp, utils.Wrap(err, "")
	}
	return resp, nil
}

func (m *AWS) upload(filePath, fileType string, onProgressFun func(int)) (string, string, error) {
	awsResp, err := m.getAwsCredentials()
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "getAwsCredentials from server failed, please check server log", err.Error(), "resp: ", awsResp)
		return "", "", utils.Wrap(err, "")
	}
	log.NewInfo("", utils.GetSelfFuncName(), "recv aws credentials", awsResp)

	newName, newType, err := m.getNewFileNameAndContentType(filePath, fileType)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "getNewFileNameAndContentType failed", err.Error(), filePath, fileType)
		return "", "", utils.Wrap(err, "")
	}
	//从接口获取参数
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(awsResp.RegionID),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsResp.AccessKeyId,
				SecretAccessKey: awsResp.SecretAccessKey,
				SessionToken:    awsResp.SessionToken,
				CanExpire:       true,
				Source:          "Open IM OSS",
			},
		}),
	)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "LoadDefaultConfig failed", err.Error(), filePath, fileType)
		return "", "", utils.Wrap(err, "")
	}
	client := s3.NewFromConfig(cfg)
	fp, err := os.Open(filePath)
	defer fp.Close()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(awsResp.Bucket),
		Key:         aws.String(newName),
		ACL:         "public-read",
		Body:        fp,
		ContentType: aws.String(newType),
	})
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "PutObject failed", err.Error(), filePath, fileType)
		return "", "", utils.Wrap(err, "")
	}
	presignedURL := "https://" + awsResp.FinalHost + "/" + newName
	return presignedURL, newName, nil
}

func (m *AWS) getNewFileNameAndContentType(filePath string, fileType string) (string, string, error) {
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

func (m *AWS) UploadImage(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "img	", onProgressFun)
}

func (m *AWS) UploadSound(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "", onProgressFun)
}

func (m *AWS) UploadFile(filePath string, onProgressFun func(int)) (string, string, error) {
	return m.upload(filePath, "", onProgressFun)
}

func (m *AWS) UploadVideo(videoPath, snapshotPath string, onProgressFun func(int)) (string, string, string, string, error) {
	videoURL, videoName, err := m.upload(videoPath, "", onProgressFun)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	snapshotURL, snapshotUUID, err := m.upload(snapshotPath, "img", onProgressFun)
	if err != nil {
		return "", "", "", "", utils.Wrap(err, "")
	}
	return snapshotURL, snapshotUUID, videoURL, videoName, nil
}
