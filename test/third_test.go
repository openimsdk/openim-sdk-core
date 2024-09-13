package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
)

type SProgress struct{}

func (s SProgress) OnProgress(current int64, size int64) {

}

func Test_UploadLog(t *testing.T) {
	tm := time.Now()
	err := open_im_sdk.UserForSDK.Third().UploadLogs(ctx, 2000, "it is ex", SProgress{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(time.Since(tm).Microseconds())

}
func Test_SDKLogs(t *testing.T) {
	open_im_sdk.UserForSDK.Third().Log(ctx, 4, "cmd/abc.go", 666, "This is a test message", "", []any{"key", "value"})
}
