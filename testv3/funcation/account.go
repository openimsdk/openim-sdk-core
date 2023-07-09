package funcation

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPB "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/log"
	"time"
)

func register(uid string) error {
	InitContext(uid)
	//ACCOUNTCHECK
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

	var rreq userPB.UserRegisterReq
	rreq.Users = []*sdkws.UserInfo{{UserID: uid}}
	for {
		err := util.ApiPost(ctx, "/auth/user_register", &rreq, nil)
		if err != nil {
			log.Error("post failed ,continue ", err.Error(), REGISTERADDR, getAccountCheckReq.CheckUserIDs)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", REGISTERADDR, getAccountCheckReq.CheckUserIDs)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		}
	}
}
