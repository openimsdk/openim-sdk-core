package v2

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
)

var _ SyncerInterface[*model_struct.LocalGroup, string] = (*Group)(nil)

func NewGroup(db db_interface.DataBase, listenerForService open_im_sdk_callback.OnListenerForService) SyncerInterface[*model_struct.LocalGroup, string] {
	return &Group{
		db:                 db,
		listenerForService: listenerForService,
	}
}

type Group struct {
	db                 db_interface.DataBase
	listenerForService open_im_sdk_callback.OnListenerForService
}

func (g *Group) OnInsert(ctx context.Context, value *model_struct.LocalGroup) error {
	return g.db.InsertGroup(ctx, value)
}

func (g *Group) OnDelete(ctx context.Context, value *model_struct.LocalGroup) error {
	return g.db.DeleteGroup(ctx, value.GroupID)
}

func (g *Group) OnUpdate(ctx context.Context, server, local *model_struct.LocalGroup) error {
	return g.db.UpdateGroup(ctx, server)
}

func (g *Group) GetID(value *model_struct.LocalGroup) string {
	return value.GroupID
}

func (g *Group) Equal(a, b *model_struct.LocalGroup) bool {
	return *a == *b
}

func (g *Group) Unchanged(ctx context.Context, value *model_struct.LocalGroup) error {
	return nil
}

func (g *Group) Insert(ctx context.Context, value *model_struct.LocalGroup) error {
	return nil
}

func (g *Group) Delete(ctx context.Context, value *model_struct.LocalGroup) error {
	return nil
}

func (g *Group) Update(ctx context.Context, server, local *model_struct.LocalGroup) error {
	return nil
}

func (g *Group) Sync(ctx context.Context, server, local []*model_struct.LocalGroup) error {
	return Sync(ctx, server, local, SyncerInterface[*model_struct.LocalGroup, string](g))
}

func Call() {
	g := NewGroup(nil, nil)

	g.Sync(context.Background(), nil, nil)

	Sync(context.Background(), nil, nil, g)

}
