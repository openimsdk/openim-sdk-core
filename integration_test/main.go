package main

import (
	"context"
	"flag"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/initialization"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/process"
	"github.com/openimsdk/tools/log"
)

func Init() {
	initialization.InitFlag()
	flag.Parse()
	initialization.SetFlagLimit()
}

func main() {
	Init()
	ctx := context.Background()
	err := process.ProV1(ctx)
	if err != nil {
		panic(err)
	}
	log.ZInfo(ctx, "start success!")
}
