package relation

import (
	"context"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/log"
)

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
