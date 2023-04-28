// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package server_api_params

type AwsStorageCredentialReq struct {
	OperationID string `json:"operationID"`
}

type AwsStorageCredentialResp struct {
	CommResp
	AccessKeyId     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	RegionID        string `json:"regionId"`
	Bucket          string `json:"bucket"`
	FinalHost       string `json:"FinalHost"`
}
