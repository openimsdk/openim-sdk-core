package db

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func (d *DataBase) InsertKey(ctx context.Context, key *model_struct.AesKey) error {
	//TODO implement me
	panic("implement me")
}

func (d *DataBase) BatchInsertKey(ctx context.Context, keys []*model_struct.AesKey) error {
	//TODO implement me
	panic("implement me")
}

func (d *DataBase) GetAesKey(ctx context.Context, sessionType int32, groupID, friendUserID string) (*model_struct.AesKey, error) {
	//TODO implement me
	panic("implement me")
}
