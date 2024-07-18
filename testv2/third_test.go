package testv2

import (
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"testing"
	"time"
)

type SProgress struct{}

func (s SProgress) OnProgress(current int64, size int64) {

}

func Test_UploadLog(t *testing.T) {
	tm := time.Now()
	err := open_im_sdk.UserForSDK.Third().UploadLogs(ctx, "测试type", "这是一个ex", SProgress{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(time.Since(tm).Microseconds())
}
