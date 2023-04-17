package chao

import "open_im_sdk/sdk_struct"

var HOST = "192.168.44.128"
var APIADDR = "http://" + HOST + ":10002"
var WSADDR = "ws://" + HOST + ":10001"

func GetConf() sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.Platform = 1
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 6
	cf.ObjectStorage = "minio"
	cf.IsCompression = true
	cf.IsExternalExtensions = true
	return cf
}
