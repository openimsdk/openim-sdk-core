package ccontext

import (
	"context"
)

type Info struct {
	UserID               string
	Token                string
	OperationID          string
	Platform             string
	ApiAddr              string
	WsAddr               string
	DataDir              string
	LogLevel             string
	ObjectStorage        string
	EncryptionKey        string
	IsCompression        bool
	IsExternalExtensions bool
}

type infoKey struct{}

func WithInfo(ctx context.Context, info *Info) context.Context {
	return context.WithValue(ctx, infoKey{}, &info)
}

func GetCtxInfo(ctx context.Context) Info {
	return *ctx.Value(infoKey{}).(*Info)
}
