package testv2

import "open_im_sdk/sdk_struct"

const (
	//APIADDR = "http://43.154.157.177:10002"
	//WSADDR  = "ws://43.154.157.177:10001"
	//UserID  = "kernaltestuid2"

	APIADDR = "http://192.168.44.128:10002"
	WSADDR  = "ws://192.168.44.128:10001"
	UserID  = "123456"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.Platform = 1
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 3
	cf.ObjectStorage = "minio"
	cf.IsCompression = true
	cf.IsExternalExtensions = true
	return cf
}
