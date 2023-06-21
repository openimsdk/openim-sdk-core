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

//go:build js && wasm
// +build js,wasm

package indexdb

import "context"

import (
	"errors"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"runtime"
	"syscall/js"
	"time"
)

//1,使用wasm原生的方式，tinygo应用于go的嵌入式领域，支持的功能有限，支持go语言的子集,甚至json序列化都无法支持
//2.函数命名遵从驼峰命名
//3.提供的sql生成语句中，关于bool值需要特殊处理，create语句的设计的到bool值的需要在创建语句中单独说明，这是因为在原有的sqlite中并不支持bool，用整数1或者0替代，gorm对其做了转换。
//4.提供的sql生成语句中，字段名是下划线方式 例如：recv_id，但是接口转换的数据json tag字段的风格是recvID，类似的所有的字段需要做个map映射
//5.任何涉及到gorm获取的是否需要返回错误，比如take和find都需要在文档上说明
//6.任何涉及到update的操作，一定要看gorm原型中实现，如果有select(*)则不需要用temp_struct中的结构体
//7.任何和接口重名的时候，db接口统一加上后缀DB
//8.任何map类型统一使用json字符串转换，文档说明

type IndexDB struct {
	LocalUsers
	LocalConversations
	*LocalChatLogs
	LocalSuperGroupChatLogs
	LocalSuperGroup
	LocalConversationUnreadMessages
	LocalGroups
	LocalGroupMember
	LocalGroupRequest
	LocalCacheMessage
	*FriendRequest
	*Black
	*Friend
	LocalChatLogReactionExtensions
	loginUserID string
}

type CallbackData struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data"`
}

const TIMEOUT = 5

var ErrTimoutFromJavaScript = errors.New("invoke javascript timeout，maybe should check  function from javascript")
var jsErr = js.Global().Get("Error")

func Exec(args ...interface{}) (output interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = utils.Wrap(errors.New(x), "")
			case error:
				err = x
			default:
				err = utils.Wrap(errors.New("unknown panic"), "")
			}
		}
	}()
	thenChannel := make(chan []js.Value)
	defer close(thenChannel)
	catchChannel := make(chan []js.Value)
	defer close(catchChannel)
	pc, _, _, _ := runtime.Caller(1)
	funcName := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	data := CallbackData{}
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js then func", "=> (main go context) "+funcName+" with respone ", args[0].String())
		thenChannel <- args
		return nil
	})
	defer thenFunc.Release()
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js catch func", "=> (main go context) "+funcName+" with respone ", args[0].String())
		catchChannel <- args
		return nil
	})
	defer catchFunc.Release()
	js.Global().Call(utils.FirstLower(funcName), args...).Call("then", thenFunc).Call("catch", catchFunc)
	select {
	case result := <-thenChannel:
		interErr := utils.JsonStringToStruct(result[0].String(), &data)
		if interErr != nil {
			err = utils.Wrap(err, "return json unmarshal err from javascript")
		}
	case catch := <-catchChannel:
		if catch[0].InstanceOf(jsErr) {
			return nil, js.Error{Value: catch[0]}
		} else {
			panic("unknown javascript exception")
		}
	case <-time.After(TIMEOUT * time.Second):
		panic(ErrTimoutFromJavaScript)
	}
	if data.ErrCode != 0 {
		return "", errors.New(data.ErrMsg)
	}
	return data.Data, err
}

func (i IndexDB) Close(ctx context.Context) error {
	_, err := Exec()
	return err
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := Exec(userID, dataDir)
	return err
}

func (i IndexDB) SetChatLogFailedStatus(ctx context.Context) {
	//msgList, err := i.GetSendingMessageList()
	//if err != nil {
	//	log.Error("", "GetSendingMessageList failed ", err.Error())
	//	return
	//}
	//for _, v := range msgList {
	//	v.Status = constant.MsgStatusSendFailed
	//	err := i.UpdateMessage(v)
	//	if err != nil {
	//		log.Error("", "UpdateMessage failed ", err.Error(), v)
	//		continue
	//	}
	//}
	//groupIDList, err := i.GetReadDiffusionGroupIDList()
	//if err != nil {
	//	log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
	//	return
	//}
	//for _, v := range groupIDList {
	//	msgList, err := i.SuperGroupGetSendingMessageList(v)
	//	if err != nil {
	//		log.Error("", "GetSendingMessageList failed ", err.Error())
	//		return
	//	}
	//	if len(msgList) > 0 {
	//		for _, v := range msgList {
	//			v.Status = constant.MsgStatusSendFailed
	//			err := i.SuperGroupUpdateMessage(v)
	//			if err != nil {
	//				log.Error("", "UpdateMessage failed ", err.Error(), v)
	//				continue
	//			}
	//		}
	//	}
	//
	//}
}
