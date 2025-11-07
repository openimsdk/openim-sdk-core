package group

import (
	"context"

	"github.com/openimsdk/protocol/sdkws"
)

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	return g.getServerJoinGroup(ctx)
}
