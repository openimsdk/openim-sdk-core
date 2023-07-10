package testv3

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/testv3/funcation"
	"testing"
)

func Test_Delay(t *testing.T) {
	var senderNum *int          // 发送者数量
	var singleSenderMsgNum *int // 单用户消息发送数量
	var intervalTime *int       // 消息发送间隔时间，ms

	senderNum = flag.Int("sn", 100, "sender num")
	singleSenderMsgNum = flag.Int("mn", 100, "single sender msg num")
	intervalTime = flag.Int("t", 1, "interval time mill second")

	flag.Parse()
	log.NewPrivateLog(funcation.LogName, funcation.LogLevel)
	log.Warn("", "reliability test start, sender num: ", *senderNum,
		" single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)

	funcation.ReliabilityTest(*singleSenderMsgNum, *intervalTime, 10, *senderNum)
}
