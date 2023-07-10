package obj_storage

import (
	"testing"
)

func TestOSS_UploadImage(t *testing.T) {
	NewOSS(nil).UploadImage("~/Desktop/images/WechatIMG21.jpeg", func(i int) {
		t.Logf("%d", i)
	})
}
