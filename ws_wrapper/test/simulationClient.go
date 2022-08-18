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

func StartSimulationJSClient(api, jssdkURL, userID string, num int, userIDList []string) {
	user := client.NewIMClient("", userID, api, jssdkURL, 1)
	var err error
	user.Token, err = user.GetToken()
	if err != nil {
		log.Println("generate token failed", userID, api, err.Error())
	}
	v := url.Values{}
	v.Set("sendID", userID)
	v.Set("token", user.Token)
	v.Set("platformID", utils.IntToString(1))
	c, _, err := websocket.DefaultDialer.Dial(jssdkURL+"?"+v.Encode(), nil)
	if err != nil {
		log.Println("dial:", err.Error(), "userID", userID, "num: ", num)
		return
	}
	lock.Lock()
	totalConnNum += 1
	log.Println("connect success", userID, "total conn num", totalConnNum)
	user.Conn = c
	lock.Unlock()
	user.WsLogout()
	user.WsLogin()
	user.GetSelfUserInfo()
	user.GetLoginStatus()
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

	go func() {
		for {
			user.SendMsg(userID)
			time.Sleep(time.Second * 1)
		}
	}()
}
