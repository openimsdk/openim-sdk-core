package msgtest

import (
	"github.com/OpenIMSDK/tools/log"
)

var (
	pressureTestAttribute PressureTestAttribute
)

func init() {
	pressureTestAttribute.InitWithFlag()
	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}
