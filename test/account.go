package test

import (
	"encoding/json"
	"net"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

func GenUid(uid int, prefix string) string {
	if getMyIP() == "" {
		log.Error("", "getMyIP() failed, exit ")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "_" + prefix + "_"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func RegisterOnlineAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "online")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)

	}
	wg.Wait()
	log.Info("", "RegisterAccounts finish ", number)
}

type GetTokenReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
}

type RegisterReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
	Name     string `json:"name"`
}

type ResToken struct {
	Data struct {
		ExpiredTime int64  `json:"expiredTime"`
		Token       string `json:"token"`
		Uid         string `json:"uid"`
	}
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

var AdminToken = ""

func init() {
	AdminToken = getToken("openIM123456")
}
func register(uid string) error {
	//ACCOUNTCHECK
	var req server_api_params.AccountCheckReq
	req.OperationID = utils.OperationIDGenerator()
	req.CheckUserIDList = []string{uid}

	var getSelfUserInfoReq server_api_params.GetSelfUserInfoReq
	getSelfUserInfoReq.OperationID = req.OperationID
	getSelfUserInfoReq.UserID = uid

	var getSelfUserInfoResp server_api_params.AccountCheckResp

	for {
		r, err := network.Post2Api(ACCOUNTCHECK, req, AdminToken)
		if err != nil {
			log.Error(req.OperationID, "post failed, continue ", err.Error(), REGISTERADDR, req)
			continue
		}
		err = json.Unmarshal(r, &getSelfUserInfoResp)
		if err != nil {
			log.Error(req.OperationID, "Unmarshal failed ", err.Error())
		}
		if getSelfUserInfoResp.ErrCode == 0 && len(getSelfUserInfoResp.ResultList) == 1 && getSelfUserInfoResp.ResultList[0].AccountStatus == "registered" {
			log.Warn(req.OperationID, "Already registered ", uid, getSelfUserInfoResp)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		} else if getSelfUserInfoResp.ErrCode == 0 && len(getSelfUserInfoResp.ResultList) == 1 && getSelfUserInfoResp.ResultList[0].AccountStatus == "unregistered" {
			log.Info(req.OperationID, "not registered ", uid, getSelfUserInfoResp.ErrCode)
			break
		} else {
			log.Error(req.OperationID, " failed, continue ", err, REGISTERADDR, req)
			continue
		}
	}

	for {
		var rreq server_api_params.UserRegisterReq
		rreq.UserID = uid
		rreq.Secret = SECRET
		rreq.UserID = uid
		rreq.Platform = 1
		rreq.OperationID = req.OperationID
		rreq.OperationID = req.OperationID
		_, err := network.Post2Api(REGISTERADDR, rreq, "")
		//if err != nil && !strings.Contains(err.Error(), "status code failed") {
		//	log.Error(req.OperationID, "post failed ,continue ", err.Error(), REGISTERADDR, req)
		//	time.Sleep(100 * time.Millisecond)
		//	continue
		//}
		if err != nil {
			log.Error(req.OperationID, "post failed ,continue ", err.Error(), REGISTERADDR, req)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info(req.OperationID, "register ok ", REGISTERADDR, req)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		}
	}
}

func getToken(uid string) string {
	url := TOKENADDR
	var req server_api_params.UserTokenReq
	req.Platform = PlatformID
	req.UserID = uid
	req.Secret = SECRET
	req.OperationID = utils.OperationIDGenerator()
	r, err := network.Post2Api(url, req, "a")
	if err != nil {
		log.Error(req.OperationID, "Post2Api failed ", err.Error(), url, req)
		return ""
	}
	var stcResp ResToken
	err = json.Unmarshal(r, &stcResp)
	if stcResp.ErrCode != 0 {
		log.Error(req.OperationID, "ErrCode failed ", stcResp.ErrCode, stcResp.ErrMsg, url, req)
		return ""
	}
	log.Info(req.OperationID, "get token: ", stcResp.Data.Token)
	return stcResp.Data.Token
}

func RunGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		} else {
			break
		}
	}
	return token
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("", "InterfaceAddrs failed ", err.Error())
		os.Exit(1)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func RegisterReliabilityUser(id int, timeStamp string) {
	userID := GenUid(id, "reliability_"+timeStamp+"_")
	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func WorkGroupRegisterReliabilityUser(id int) {
	userID := GenUid(id, "workgroup")
	//	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	log.Info("", "WorkGroupRegisterReliabilityUser userID: ", userID, "token: ", token)
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func RegisterPressUser(id int) {
	userID := GenUid(id, "press")
	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func GetGroupMemberNum(groupID string) uint32 {
	var req server_api_params.GetGroupInfoReq
	req.OperationID = utils.OperationIDGenerator()
	req.GroupIDList = []string{groupID}

	var groupInfoList []*server_api_params.GroupInfo

	r, err := network.Post2Api(GETGROUPSINFOROUTER, req, AdminToken)
	if err != nil {
		log.Error("", "post failed ", GETGROUPSINFOROUTER, req)
		return 0
	}
	err = common.CheckErrAndResp(nil, r, &groupInfoList, nil)
	if err != nil {
		log.Error("", "CheckErrAndResp failed ", err.Error(), string(r))
		return 0
	}
	log.Warn("", "group info", groupInfoList)
	return groupInfoList[0].MemberCount
}
