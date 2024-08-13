package decorator

import (
	"context"
	"fmt"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
	"time"
)

// FuncLog is a log print decorator.
// Correct usage is: defer decorator.FuncLog(ctx)()
func FuncLog(ctx context.Context) func() {
	return FuncLogSkip(ctx, 1)
}

// FuncLogSkip is a log print decorator. The argument skip is the number of stack frames
// to ascend.
// e.g.
//
//	func FuncName(ctx context.Context){
//		   middleFunc(ctx)
//	}
//
//	func middleFunc(ctx context.Context){
//	    FuncLogSkip(ctx, 1)
//	    // ...
//	}
//
// the funcName is `FuncName`
func FuncLogSkip(ctx context.Context, skip int) func() {
	t := time.Now()
	funcName := stringutil.GetFuncName(skip + 1) // +1 is FuncLogSkip
	log.ZInfo(ctx, fmt.Sprintf("%s begin", funcName))
	fmt.Println(fmt.Sprintf("%s begin", funcName))
	return func() {
		log.ZInfo(ctx, fmt.Sprintf("%s end", funcName), "time consuming", time.Since(t))
		fmt.Println(fmt.Sprintf("%s end. Time consuming: %v", funcName, time.Since(t)))
	}
}
