// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package server_api_params

type OSSCredentialReq struct {
	OperationID string `json:"operationID"`
	Filename    string `json:"filename"`
	FileType    string `json:"file_type"`
}

type OSSCredentialRespData struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Token           string `json:"token"`
	Bucket          string `json:"bucket"`
	FinalHost       string `json:"final_host"`
}

type OSSCredentialResp struct {
	CommResp
	OssData OSSCredentialRespData  `json:"-"`
	Data    map[string]interface{} `json:"data"`
}
