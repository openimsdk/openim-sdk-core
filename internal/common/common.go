// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package common

import (
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

func UnmarshalTips(msg *sdkws.MsgData, detail proto.Message) error {
	var tips sdkws.TipsComm
	if err := proto.Unmarshal(msg.Content, &tips); err != nil {
		return utils.Wrap(err, "")
	}
	if err := proto.Unmarshal(tips.Detail, detail); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}
