// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func Test_CreateTextMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateTextMessage(ctx, "textMsg")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateAdvancedTextMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateAdvancedTextMessage(ctx, "textAdMsg", []*sdk_struct.MessageEntity{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateTextAtMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateTextAtMessage(ctx, "textATtsg", []string{}, []*sdk_struct.AtInfo{}, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateQuoteMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateQuoteMessage(ctx, "textATtsg", &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateAdvancedQuoteMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateAdvancedQuoteMessage(ctx, "textATtsg", &sdk_struct.MsgStruct{}, []*sdk_struct.MessageEntity{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessageFromFullPath(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateVideoMessageFromFullPath(ctx, ".\\test.png", "mp4", 10, ".\\test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateCardMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateCardMessage(ctx, &sdk_struct.CardElem{
		UserID:   "123456",
		Nickname: "testname",
		FaceURL:  "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateImageMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateImageMessage(ctx, ".\\test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateImageMessageByURL(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateImageMessageByURL(ctx, "",
		sdk_struct.PictureBaseInfo{}, sdk_struct.PictureBaseInfo{}, sdk_struct.PictureBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateSoundMessageByURL(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateSoundMessageByURL(ctx, &sdk_struct.SoundBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateSoundMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateSoundMessage(ctx, ".\\test.png", 20)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateVideoMessage(ctx, ".\\test.png", "mp4", 10, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessageByURL(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateVideoMessageByURL(ctx, sdk_struct.VideoBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFileMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateFileMessage(ctx, ".\\test.png", "png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFileMessageByURL(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateFileMessageByURL(ctx, sdk_struct.FileBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateLocationMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateLocationMessage(ctx, "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateCustomMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateCustomMessage(ctx, "", "", "")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateMergerMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateMergerMessage(ctx, []*sdk_struct.MsgStruct{}, "title", []string{"summary"})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFaceMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateFaceMessage(ctx, 0, "www.faceURL.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateForwardMessage(t *testing.T) {
	message, err := open_im_sdk.IMUserContext.Conversation().CreateForwardMessage(ctx, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}
