package common

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"open_im_sdk/pkg/common"
)

type SendMsgCallBack interface {
	common.Base
	OnProgress(progress int)
}

type ObjectStorage interface {
	UploadImage(filePath string, callback SendMsgCallBack, operationID string) (string, string)
	UploadSound(filePath string, callback SendMsgCallBack, operationID string) (string, string)
	UploadFile(filePath string, callback SendMsgCallBack, operationID string) (string, string)
	UploadVideo(videoPath, snapshotPath string, callback SendMsgCallBack, operationID string) (string, string, string, string)
}
