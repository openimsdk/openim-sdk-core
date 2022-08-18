package test

import (
	"log"
	"net/url"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/ws_wrapper/test/client"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var totalConnNum int
var lock sync.Mutex

func StartSimulationJSClient(api, jssdkURL, userID string, i int, userIDList []string) {
	user := client.NewIMClient("", userID, api, jssdkURL, 5)
	var err error
	user.Token, err = user.GetToken()
	if err != nil {
		log.Println("generate token failed", userID, api, err.Error())
	}
	v := url.Values{}
	v.Set("sendID", userID)
	v.Set("token", user.Token)
	v.Set("platformID", utils.IntToString(5))
	c, _, err := websocket.DefaultDialer.Dial(jssdkURL+"?"+v.Encode(), nil)
	if err != nil {
		log.Println("dial:", err.Error(), "userID", userID, "i: ", i)
		return
	}
	lock.Lock()
	totalConnNum += 1
	log.Println("connect success", userID, "total conn num", totalConnNum)
	lock.Unlock()
	user.Conn = c
	// user.WsLogout()
	user.WsLogin()
	time.Sleep(time.Second * 2)

	// 模拟同步
	go func() {
		user.GetSelfUserInfo()
		user.GetAllConversationList()
		user.GetBlackList()
		user.GetFriendList()
		user.GetRecvFriendApplicationList()
		user.GetRecvGroupApplicationList()
		user.GetSendFriendApplicationList()
		user.GetSendGroupApplicationList()
	}()

	// 模拟监听回调
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err, "error an connet failed", userID)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// 模拟给随机用户发消息
	go func() {
		for {
			user.SendMsg(userID)
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			user.GetLoginStatus()
			time.Sleep(time.Second * 10)
		}
	}()
}
