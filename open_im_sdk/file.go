package open_im_sdk

import (
	"open_im_sdk/internal/file"
	"open_im_sdk/open_im_sdk_callback"
)

func PutFile(callback open_im_sdk_callback.Base, operationID string, req string, progress open_im_sdk_callback.PutFileCallback) {
	call(callback, operationID, UserForSDK.File().PutFile, req, file.PutFileCallback(progress))
}
