// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package obj_storage

import (
	"testing"
)

func TestOSS_UploadImage(t *testing.T) {
	NewOSS(nil).UploadImage("~/Desktop/images/WechatIMG21.jpeg", func(i int) {
		t.Logf("%d", i)
	})
}
