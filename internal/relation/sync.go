package relation

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (r *Relation) SyncBothFriendRequest(ctx context.Context, fromUserID, toUserID string) error {
	if toUserID == r.loginUserID {
		if !r.requestRecvSyncerLock.TryLock() {
			return nil
		}
		defer r.requestRecvSyncerLock.Unlock()
	} else {
		if !r.requestSendSyncerLock.TryLock() {
			return nil
		}
		defer r.requestSendSyncerLock.Unlock()
	}
	req := &relation.GetDesignatedFriendsApplyReq{FromUserID: fromUserID, ToUserID: toUserID}
	friendRequests, err := r.getDesignatedFriendsApply(ctx, req)
	if err != nil {
		return err
	}
	localData, err := r.db.GetBothFriendReq(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}
	if toUserID == r.loginUserID {
		return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, friendRequests), localData, nil)
	} else if fromUserID == r.loginUserID {
		return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, friendRequests), localData, nil)
	}
	return nil
}

// SyncAllSelfFriendApplication send
func (r *Relation) SyncAllSelfFriendApplication(ctx context.Context) error {
	if !r.requestSendSyncerLock.TryLock() {
		return nil
	}
	defer r.requestSendSyncerLock.Unlock()
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	requests, err := r.getSelfFriendApplicationList(ctx, req)
	if err != nil {
		return err
	}
	localData, err := r.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

func (r *Relation) SyncAllSelfFriendApplicationWithoutNotice(ctx context.Context) error {
	if !r.requestSendSyncerLock.TryLock() {
		return nil
	}
	defer r.requestSendSyncerLock.Unlock()
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	requests, err := r.getSelfFriendApplicationList(ctx, req)
	if err != nil {
		return err
	}
	localData, err := r.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

// SyncAllFriendApplication recv
func (r *Relation) SyncAllFriendApplication(ctx context.Context) error {
	if !r.requestRecvSyncerLock.TryLock() {
		return nil
	}
	defer r.requestRecvSyncerLock.Unlock()
	req := &relation.GetPaginationFriendsApplyToReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	requests, err := r.getFriendApplicationList(ctx, req)
	if err != nil {
		return err
	}
	localData, err := r.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}
func (r *Relation) SyncAllFriendApplicationWithoutNotice(ctx context.Context) error {
	if !r.requestRecvSyncerLock.TryLock() {
		return nil
	}
	defer r.requestRecvSyncerLock.Unlock()
	req := &relation.GetPaginationFriendsApplyToReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	requests, err := r.getFriendApplicationList(ctx, req)
	if err != nil {
		return err
	}
	localData, err := r.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

func (r *Relation) SyncAllFriendList(ctx context.Context) error {
	t := time.Now()
	defer func(start time.Time) {
		elapsed := time.Since(start).Milliseconds()
		log.ZDebug(ctx, "SyncAllFriendList fn call end", "cost time", fmt.Sprintf("%d ms", elapsed))

	}(t)
	return r.IncrSyncFriends(ctx)
}

func (r *Relation) SyncAllBlackList(ctx context.Context) error {
	serverData, err := r.getBlackList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return r.blackSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}

func (r *Relation) SyncAllBlackListWithoutNotice(ctx context.Context) error {
	serverData, err := r.getBlackList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return r.blackSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil, false, true)
}

func (r *Relation) GetDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	return r.getDesignatedFriends(ctx, friendIDs)
}
