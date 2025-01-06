package db

import (
	"context"
	"testing"
)

func TestGetLatestValidateServerMessage(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "../../", 6)
	if err != nil {
		return
	}

	message, err := db.GetLatestValidServerMessage(ctx, "sg_93606743", 1710774468519, true)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("message", message)
}
