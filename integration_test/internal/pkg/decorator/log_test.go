package decorator

import (
	"context"
	"testing"
)

func TestFuncLog(t *testing.T) {
	FuncName(context.Background())
}

func FuncName(ctx context.Context) {
	middleFunc(ctx)
}

func middleFunc(ctx context.Context) {
	FuncLogSkip(ctx, 1)
	//...
}
