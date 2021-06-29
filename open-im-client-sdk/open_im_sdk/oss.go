package open_im_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"time"
)

func tencentOssCredentials() (*paramsTencentOssCredentialResp, error) {
	resp, err := post2Api(tencentCloudStorageCredentialRouter, paramsTencentOssCredentialReq{OperationID: operationIDGenerator()}, token)
	if err != nil {
		return nil, err
	}

	var ossResp paramsTencentOssCredentialResp
	_ = json.Unmarshal(resp, &ossResp)

	if ossResp.ErrCode != 0 {
		return nil, errors.New(ossResp.ErrMsg)
	}

	return &ossResp, nil
}

func uploadImage(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := tencentOssCredentials()
	if err != nil {
		sdkLog("tencentOssCredentials", err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, err := url.Parse(dir)
	if err != nil {
		sdkLog("Parse", err.Error())
		return "", "", err
	}
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {

		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: contentType,
				Listener:    lis,
			},
		}
		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog("file:", filePath, err.Error())
			return "", "", err
		}

		targetFileUrl := dir + "/" + newName
		return targetFileUrl, newName, nil
	}

	return "", "", errors.New("client == nil")
}

func uploadSound(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := tencentOssCredentials()
	if err != nil {
		sdkLog(err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {

		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog("PutFromFile", err.Error())
			return "", "", err
		}

		targetFile := dir + "/" + newName
		return targetFile, newName, nil
	}
	sdkLog("client == nil")
	return "", "", errors.New("client == nil")
}

func uploadFile(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := tencentOssCredentials()
	if err != nil {
		sdkLog(err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {

		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog(err.Error())
			return "", "", err
		}

		targetFile := dir + "/" + newName
		return targetFile, newName, nil
	}

	return "", "", errors.New("client == nil")
}

func uploadVideo(videoPath, snapshotPath string, back SendMsgCallBack) (string, string, string, string, error) {
	ossResp, err := tencentOssCredentials()
	if err != nil {
		sdkLog(err.Error())
		return "", "", "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {
		var newNameSnapshot, targetSnapshot string
		if len(snapshotPath) > 0 {
			//-----first------
			suffix := path.Ext(snapshotPath)
			if len(suffix) == 0 {
				return "", "", "", "", errors.New("file fail")
			}
			newNameSnapshot := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
			contentTypeSnapshot := "image/" + suffix[1:]

			opt1 := &cos.ObjectPutOptions{
				ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
					ContentType: contentTypeSnapshot,
				},
			}

			_, err := client.Object.PutFromFile(context.Background(), newNameSnapshot, snapshotPath, opt1)
			if err != nil {
				sdkLog(err.Error())
				return "", "", "", "", err
			}

			targetSnapshot = dir + "/" + newNameSnapshot
		}

		//-----second------
		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(videoPath)
		if len(suffix) == 0 {
			return "", "", "", "", errors.New("file fail")
		}
		newNameVideo := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt2 := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err = client.Object.PutFromFile(context.Background(), newNameVideo, videoPath, opt2)
		if err != nil {
			sdkLog(err.Error())
			return "", "", "", "", err
		}

		targetVideo := dir + "/" + newNameVideo
		return targetSnapshot, newNameSnapshot, targetVideo, newNameVideo, nil
	}

	return "", "", "", "", errors.New("client == nil")
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
		log(fmt.Sprintf("\r[ConsumedBytes/TotalBytes: %d/%d, %d%%]", event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes))

	case cos.ProgressFailedEvent:
		log(fmt.Sprintf("\nTransfer Failed: %v", event.Err))
	}
}
