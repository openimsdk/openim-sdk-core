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

package ccontext

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/tools/mcontext"
)

const (
	Callback = "callback"
)

type GlobalConfig struct {
	UserID string
	Token  string

	sdk_struct.IMConfig
}

type ContextInfo interface {
	UserID() string
	Token() string
	PlatformID() int32
	ApiAddr() string
	WsAddr() string
	DataDir() string
	LogLevel() uint32
	OperationID() string
	IsExternalExtensions() bool
}

func Info(ctx context.Context) ContextInfo {
	conf := ctx.Value(GlobalConfigKey{}).(*GlobalConfig)
	return &info{
		conf: conf,
		ctx:  ctx,
	}
}

func WithInfo(ctx context.Context, conf *GlobalConfig) context.Context {
	return context.WithValue(ctx, GlobalConfigKey{}, conf)
}

func WithOperationID(ctx context.Context, operationID string) context.Context {
	return mcontext.SetOperationID(ctx, operationID)
}
func WithSendMessageCallback(ctx context.Context, callback open_im_sdk_callback.SendMsgCallBack) context.Context {
	return context.WithValue(ctx, Callback, callback)
}

type GlobalConfigKey struct{}

type info struct {
	conf *GlobalConfig
	ctx  context.Context
}

func (i *info) UserID() string {
	return i.conf.UserID
}

func (i *info) Token() string {
	return i.conf.Token
}

func (i *info) PlatformID() int32 {
	return i.conf.PlatformID
}

func (i *info) ApiAddr() string {
	return i.conf.ApiAddr
}

func (i *info) WsAddr() string {
	return i.conf.WsAddr
}

func (i *info) DataDir() string {
	return i.conf.DataDir
}

func (i *info) LogLevel() uint32 {
	return i.conf.LogLevel
}

func (i *info) OperationID() string {
	return mcontext.GetOperationID(i.ctx)
}

func (i *info) IsExternalExtensions() bool {
	return i.conf.IsExternalExtensions
}
