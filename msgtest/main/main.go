package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/module"
)

func init() {
	_ = runtime.GOMAXPROCS(7)
	InitWithFlag()
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "./", 2, 24); err != nil {
		panic(err)
	}
}

var (
	totalOnlineUserNum    int     // 总在线用户数
	randomSender          int     // 随机发送者数
	randomReceiver        int     // 随机接收者数
	samplingRate          float64 // 抽样率
	NotFriendMsgSenderNum int     // 非好友消息发送者数
	groupMsgSenderNum     int     // 群消息发送者数
	msgSenderNumEvreyUser int     // 每个用户的消息数
	fastenedUserNum       int     // 固定用户数
	start                 int
	end                   int
	count                 int
	sendInterval          int
    onlineUsersOnly bool


	//recvMsgUserNum int // 消息接收者数, 抽样账号
	isRegisterUser bool // 是否注册用户
)

func InitWithFlag() {
	flag.IntVar(&totalOnlineUserNum, "o", 20000, "total online user num")
	flag.IntVar(&randomSender, "rs", 100, "random sender num")
	flag.IntVar(&randomReceiver, "rr", 100, "random receiver num")
	flag.IntVar(&start, "s", 0, "start user")
	flag.IntVar(&end, "e", 0, "end user")
	flag.Float64Var(&samplingRate, "f", 0.1, "sampling rate")
	flag.IntVar(&count, "c", 200, "number of messages per user")
	flag.IntVar(&sendInterval, "i", 1000, "send message interval per user(milliseconds)")
	flag.IntVar(&NotFriendMsgSenderNum, "n", 100, "not friend msg sender num")
	flag.IntVar(&groupMsgSenderNum, "g", 100, "group msg sender num")
	flag.IntVar(&msgSenderNumEvreyUser, "m", 100, "msg sender num evrey user")

	flag.BoolVar(&isRegisterUser, "r", false, "register user to IM system")
    flag.BoolVar(&onlineUsersOnly, "u", false, "consider only online users")

}

func PrintQPS() {
	for {

		log.ZError(context.Background(), "QPS", nil, "qps", module.GetQPS())
		time.Sleep(time.Second * 1)
	}
}

func main() {
	flag.Parse()
	fmt.Print("start", totalOnlineUserNum, count, sendInterval, isRegisterUser,onlineUsersOnly)
	ctx := context.Background()
	p := module.NewPressureTester()
	f, r, err := p.SelectSample(totalOnlineUserNum, 0.01)
	//f, r, err := p.SelectSample2(totalOnlineUserNum, 0.01)
	if err != nil {
		log.ZError(ctx, "Sample UserID failed", err)
		return
	}
	log.ZDebug(ctx, "Sample UserID", "sampleUserLength", len(r), "sampleUserID", r, "length", len(f))
	time.Sleep(10 * time.Second)
	//
	if isRegisterUser {
		if err := p.RegisterUsers(f, nil, nil); err != nil {
			log.ZError(ctx, "Sample UserID failed", err)
			return
		}
	}
	if start != 0 {
		f = p.SelectStartAndEnd(start, end)
	}
	//go PrintQPS()
	// init users
	p.InitUserConns(f)
	log.ZWarn(ctx, "all user init connect to server success,start send message", nil, "count", count)
	if onlineUsersOnly {
	log.ZWarn(ctx, "Blocking the process...", nil)
		// Create a channel to receive operating system interrupt signals
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

		// Block the process until an interrupt signal is received
		<-signalChannel
	log.ZWarn(ctx, "Received interrupt signal. Exiting...", nil)
		return
	}
	time.Sleep(10 * time.Second)
	p.SendSingleMessages2(f, p.Shuffle(f, randomSender), randomReceiver, count, time.Millisecond*time.Duration(sendInterval))
	log.ZWarn(ctx, "send over", nil, "num", p.GetSendNum())
	//p.SendSingleMessagesTo(f, 20000, time.Millisecond*1)
	//p.SendMessages("fastened_user_prefix_testv3new_0", "fastened_user_prefix_testv3new_1", 100000)
	time.Sleep(1 * time.Minute)
	p.CheckMsg(ctx)

	time.Sleep(time.Hour * 60)

}
