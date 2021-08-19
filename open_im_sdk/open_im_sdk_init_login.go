package open_im_sdk

import "encoding/json"

type IMSDKListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(ErrCode int, ErrMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
	OnSelfInfoUpdated(userInfo string)
}

func InitSDK(config string, cb IMSDKListener) bool {
	return SdkInitManager.initSDK(config, cb)

}
func UnInitSDK() {
	SdkInitManager.unInitSDK()
}

func Login(uid, tk string, callback Base) {
	if callback == nil {
		sdkLog("callback is null")
		return
	}
	go func() {
		SdkInitManager.login(uid, tk, callback)
	}()
}

func ForceReConn() {
	if SdkInitManager.conn != nil {
		SdkInitManager.conn.Close()
	}
}

func Logout(callback Base) {
	SdkInitManager.logout(callback)
}

func GetLoginStatus() int {
	return SdkInitManager.getLoginStatus()
}

func GetLoginUser() string {
	return LoginUid
}

func ForceSyncLoginUerInfo() {
	SdkInitManager.syncLoginUserInfo()
}

type Base interface {
	OnError(errCode int, errMsg string)
	OnSuccess(data string)
}

func TencentOssCredentials(cb Base) {
	resp, err := tencentOssCredentials()
	if err != nil {
		cb.OnError(-1, err.Error())
		return
	}
	bResp, err := json.Marshal(resp)
	if err != nil {
		cb.OnError(-1, err.Error())
		return
	}
	cb.OnSuccess(string(bResp))
}
