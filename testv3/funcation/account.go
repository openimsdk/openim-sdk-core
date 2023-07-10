package funcation

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/log"
	"time"
)

func Register(uid, nickname, faceurl string) (bool, error) {
	InitContext(uid)
	err := checkUserAccount(uid)
	return registerUserAccount(uid, nickname, faceurl), err
}

func checkUserAccount(uid string) error {
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{uid}
	for {
		err := util.ApiPost(ctx, "/user/account_check", &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.Warn(getAccountCheckReq.CheckUserIDs[0], "Already registered ", uid, getAccountCheckResp.Results)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.Info(getAccountCheckReq.CheckUserIDs[0], "not registered ", uid, getAccountCheckResp.Results)
			break
		} else {
			log.Error(getAccountCheckReq.CheckUserIDs[0], " failed, continue ", err, REGISTERADDR, getAccountCheckReq.CheckUserIDs)
			continue
		}
	}
	return nil
}

func registerUserAccount(uid, nickname, faceurl string) bool {
	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: uid, Nickname: nickname, FaceURL: faceurl}}
	req.Secret = Secret
	for {
		err := util.ApiPost(ctx, "/user/user_register", &req, nil)
		if err != nil {
			log.Error("post failed ,continue ", err.Error(), REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", REGISTERADDR)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return true
		}
	}
	return false
}
