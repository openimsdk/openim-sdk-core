package decorator

import (
	"context"
	"fmt"
	"testing"
)

func TestFuncLog(t *testing.T) {
	FuncName(context.Background())
}

func FuncName(ctx context.Context) {
	middleFunc(ctx)
}

func middleFunc(ctx context.Context) {
	defer FuncLogSkip(ctx, 1)()
	//...
	fmt.Println("middle func")
}
