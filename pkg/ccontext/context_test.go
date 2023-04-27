// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package ccontext

import (
	"context"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	conf := GlobalConfig{
		UserID: "uid123",
	}
	ctx = WithInfo(ctx, &conf)
	operationID := "opid123"
	ctx = WithOperationID(ctx, operationID)

	fmt.Println("UserID:", Info(ctx).UserID())
	fmt.Println("OperationID:", Info(ctx).OperationID())
	if Info(ctx).UserID() != conf.UserID {
		t.Fatal("UserID not match")
	}
	if Info(ctx).OperationID() != operationID {
		t.Fatal("OperationID not match")
	}
}
