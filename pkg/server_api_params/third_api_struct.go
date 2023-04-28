// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package server_api_params

import sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"

type TencentCloudStorageCredentialReq struct {
	OperationID string `json:"operationID"`
}

type TencentCloudStorageCredentialRespData struct {
	*sts.CredentialResult
	Region string `json:"region"`
	Bucket string `json:"bucket"`
}

type TencentCloudStorageCredentialResp struct {
	CommResp
	CosData TencentCloudStorageCredentialRespData `json:"-"`
	Data    map[string]interface{}                `json:"data"`
}
type FcmUpdateTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Platform    int    `json:"platform" binding:"required,min=1,max=2"` // only for ios + android
	FcmToken    string `json:"fcmToken" binding:"required"`
}

type FcmUpdateTokenResp struct {
	CommResp
}

type SetAppBadgeReq struct {
	OperationID    string `json:"operationID" binding:"required"`
	FromUserID     string `json:"fromUserID" binding:"required"`
	AppUnreadCount int32  `json:"appUnreadCount" binding:"required"`
}

type SetAppBadgeResp struct {
	CommResp
}
