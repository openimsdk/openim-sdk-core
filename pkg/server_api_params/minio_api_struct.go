package server_api_params

type MinioStorageCredentialReq struct {
	OperationID string `json:"operationID"`
}

type MinioStorageCredentialResp struct {
	CommResp
	SecretAccessKey string `json:"secretAccessKey"`
	AccessKeyID     string `json:"accessKeyID"`
	SessionToken    string `json:"sessionToken"`
	SignerType      int    `json:"signerType"`
	BucketName      string `json:"bucketName"`
	StsEndpointURL  string `json:"stsEndpointURL"`
	StorageTime     int    `json:"storageTime"`
}
