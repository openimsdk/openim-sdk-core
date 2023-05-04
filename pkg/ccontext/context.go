// Copyright © 2023 OpenIM SDK.
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
)

type GlobalConfig struct {
	UserID               string
	Token                string
	Platform             int32
	ApiAddr              string
	WsAddr               string
	DataDir              string
	LogLevel             uint32
	EncryptionKey        string
	IsCompression        bool
	IsExternalExtensions bool
}

type ContextInfo interface {
	UserID() string
	Token() string
	Platform() int32
	ApiAddr() string
	WsAddr() string
	DataDir() string
	LogLevel() uint32
	EncryptionKey() string
	OperationID() string
	IsCompression() bool
	IsExternalExtensions() bool
}

func Info(ctx context.Context) ContextInfo {
	conf := ctx.Value(globalConfigKey{}).(*GlobalConfig)
	return &info{
		conf: conf,
		ctx:  ctx,
	}
}

func WithInfo(ctx context.Context, conf *GlobalConfig) context.Context {
	return context.WithValue(ctx, globalConfigKey{}, conf)
}

func WithOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, operationIDKey, operationID)
}

type globalConfigKey struct{}

const operationIDKey = "operationID" // 兼容服务端

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

func (i *info) Platform() int32 {
	return i.conf.Platform
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

func (i *info) EncryptionKey() string {
	return i.conf.EncryptionKey
}

func (i *info) OperationID() string {
	return i.ctx.Value(operationIDKey).(string)
}

func (i *info) IsCompression() bool {
	return i.conf.IsCompression
}

func (i *info) IsExternalExtensions() bool {
	return i.conf.IsExternalExtensions
}
