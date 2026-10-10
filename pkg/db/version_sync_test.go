package db

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func Test_GetVersionSync(t *testing.T) {

	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}

	_, err = db.GetVersionSync(ctx, "local_group_entities_version", "1076204769")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	// log.ZDebug(ctx, "info is ", "info", info)
}

func Test_SetVersionSync(t *testing.T) {

	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}

	newVersionSync := &model_struct.LocalVersionSync{
		Table:      "local_group_entities_version",
		EntityID:   "1076204769",
		VersionID:  "667aabe3417b67f0f0d3cdee",
		Version:    1076204769,
		CreateTime: 0,
		UIDList:    []string{"8879166186", "1695766238", "2882899447", "5292156665"},
	}

	err = db.SetVersionSync(ctx, newVersionSync)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	// log.ZDebug(ctx, "info is ", "info", info)
}

func Test_DeleteVersionSync(t *testing.T) {

	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}

	err = db.DeleteVersionSync(ctx, "local_group_entities_version", "3378458183")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	// log.ZDebug(ctx, "info is ", "info", info)
}
