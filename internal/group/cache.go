package group

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/datafetcher"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) buildGroupMemberKey(groupID string, userID string) string {
	return groupID + ":" + userID
}

func (g *Group) GetGroupMembersInfoFunc(ctx context.Context, groupID string, userIDs []string,
	fetchFunc func(ctx context.Context, missingKeys []string) ([]*model_struct.LocalGroupMember, error),
) (map[string]*model_struct.LocalGroupMember, error) {
	var (
		res         = make(map[string]*model_struct.LocalGroupMember)
		missingKeys []string
	)

	for _, userID := range userIDs {
		key := g.buildGroupMemberKey(groupID, userID)
		if member, ok := g.groupMemberCache.Load(key); ok {
			res[userID] = member
		} else {
			missingKeys = append(missingKeys, userIDs...)
		}
	}

	log.ZDebug(ctx, "GetGroupMembersInfoFunc fetch", "missingKeys", missingKeys)
	fetchData, err := fetchFunc(ctx, missingKeys)
	if err != nil {
		return nil, err
	}

	for i, data := range fetchData {
		key := g.buildGroupMemberKey(groupID, data.UserID)
		res[data.UserID] = fetchData[i]
		g.groupMemberCache.Store(key, fetchData[i])
	}

	return res, nil
}

func (g *Group) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) (map[string]*model_struct.LocalGroupMember, error) {
	res := make(map[string]*model_struct.LocalGroupMember)
	if len(userIDs) == 0 {
		return res, nil
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(member *model_struct.LocalGroupMember) string {
			return member.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return nil
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			var (
				localData []*model_struct.LocalGroupMember
				needDB    []string
			)

			for _, userID := range userIDs {
				key := g.buildGroupMemberKey(groupID, userID)
				if member, ok := g.groupMemberCache.Load(key); ok {
					localData = append(localData, member)
				} else {
					needDB = append(needDB, userID)
				}
			}

			if len(needDB) == 0 {
				return localData, false, nil
			}

			dbData, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, needDB)
			if err != nil {
				return nil, false, err
			}

			for _, member := range dbData {
				g.groupMemberCache.Store(g.buildGroupMemberKey(groupID, member.UserID), member)
				localData = append(localData, member)
			}

			if len(dbData) == len(needDB) {
				return localData, false, nil
			}

			return localData, true, nil
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			if len(userIDs) == 0 {
				return nil, nil
			}
			queryData, err := g.getDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			converted := datautil.Batch(ServerGroupMemberToLocalGroupMember, queryData)
			for _, member := range converted {
				g.groupMemberCache.Store(g.buildGroupMemberKey(groupID, member.UserID), member)
			}
			return converted, nil
		},
	)

	fetchData, err := dataFetcher.FetchMissingAndCombineLocal(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for _, member := range fetchData {
		res[member.UserID] = member
	}
	return res, nil
}
