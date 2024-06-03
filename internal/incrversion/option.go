package incrversion

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"time"
)

type Option[V, R any] struct {
	Ctx           context.Context
	DB            db_interface.VersionSyncModel
	Key           func(V) string
	SyncKey       func() string
	Local         func() ([]V, error)
	ServerVersion func() R
	Server        func(version *model_struct.LocalVersionSync) (R, error)
	Full          func(resp R) bool
	Version       func(resp R) (string, uint64)
	DeleteIDs     func(resp R) []string
	Changes       func(resp R) []V
	Syncer        func(server, local []V) error
}

func (o *Option[V, R]) getVersionInfo() *model_struct.LocalVersionSync {
	key := o.SyncKey()
	versionInfo, err := o.DB.GetVersionSync(o.Ctx, key)
	if err != nil {
		log.ZInfo(o.Ctx, "get version info", "error", err)
		return &model_struct.LocalVersionSync{
			Key: key,
		}
	}
	return versionInfo
}

func (o *Option[V, R]) Sync() error {
	var resp R
	if o.ServerVersion == nil {
		version := o.getVersionInfo()
		var err error
		resp, err = o.Server(version)
		if err != nil {
			return err
		}
	} else {
		resp = o.ServerVersion()
	}
	delIDs := o.DeleteIDs(resp)
	changes := o.Changes(resp)
	updateVersionInfo := func() error {
		lvs := &model_struct.LocalVersionSync{
			Key:        o.SyncKey(),
			CreateTime: time.Now().UnixMilli(),
		}
		lvs.VersionID, lvs.Version = o.Version(resp)
		return o.DB.SetVersionSync(o.Ctx, lvs)
	}
	if len(delIDs)+len(changes) == 0 {
		return updateVersionInfo()
	}
	local, err := o.Local()
	if err != nil {
		return err
	}
	var server []V
	if o.Full(resp) {
		server = changes
	} else {
		kv := datautil.SliceToMapAny(local, func(v V) (string, V) {
			return o.Key(v), v
		})
		for i, change := range changes {
			key := o.Key(change)
			kv[key] = changes[i]
		}
		for _, id := range delIDs {
			delete(kv, id)
		}
		server = datautil.Values(kv)
	}
	if err := o.Syncer(server, local); err != nil {
		return err
	}
	return updateVersionInfo()
}
