package full

import (
	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/relation"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
)

type Full struct {
	user     *user.User
	relation *relation.Relation
	group    *group.Group
	ch       chan common.Cmd2Value
	db       db_interface.DataBase
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, relation *relation.Relation, group *group.Group, ch chan common.Cmd2Value,
	db db_interface.DataBase) *Full {
	return &Full{user: user, relation: relation, group: group, ch: ch, db: db}
}
