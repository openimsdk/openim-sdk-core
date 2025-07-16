package syncer

import (
	"context"
	"reflect"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type VersionSynchronizer[V, R any] struct {
	Ctx                context.Context
	DB                 db_interface.VersionSyncModel
	TableName          string
	EntityID           string
	Key                func(V) string
	Local              func() ([]V, error)
	ServerVersion      func() R
	Server             func(version *model_struct.LocalVersionSync) (R, error)
	Full               func(resp R) bool
	Version            func(resp R) (string, uint64)
	Delete             func(resp R) []string
	Update             func(resp R) []V
	Insert             func(resp R) []V
	ExtraData          func(resp R) any
	ExtraDataProcessor func(ctx context.Context, data any) error
	Syncer             func(server, local []V) error
	FullSyncer         func(ctx context.Context) error
	FullID             func(ctx context.Context) ([]string, error)
	IDOrderChanged     func(resp R) bool
}

func (o *VersionSynchronizer[V, R]) getVersionInfo() (*model_struct.LocalVersionSync, error) {
	versionInfo, err := o.DB.GetVersionSync(o.Ctx, o.TableName, o.EntityID)
	if err != nil && !errs.ErrRecordNotFound.Is(errs.Unwrap(err)) {
		log.ZWarn(o.Ctx, "get version info", err)
		return nil, err

	}
	return versionInfo, nil
}

func (o *VersionSynchronizer[V, R]) updateVersionInfo(lvs *model_struct.LocalVersionSync, resp R) error {
	lvs.Table = o.TableName
	lvs.EntityID = o.EntityID
	lvs.VersionID, lvs.Version = o.Version(resp)
	return o.DB.SetVersionSync(o.Ctx, lvs)
}
func judgeInterfaceIsNil(data any) bool {
	return reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()
}

func (o *VersionSynchronizer[V, R]) IncrementalSync() error {
	var lvs *model_struct.LocalVersionSync
	var resp R
	var extraData any
	if o.ServerVersion == nil {
		var err error
		lvs, err = o.getVersionInfo()
		if err != nil {
			return err
		}
		resp, err = o.Server(lvs)
		if err != nil {
			return err
		}
	} else {
		var err error
		lvs, err = o.getVersionInfo()
		if err != nil {
			return err
		}
		resp = o.ServerVersion()
	}

	delIDs := o.Delete(resp)
	changes := o.Update(resp)
	insert := o.Insert(resp)
	if o.ExtraData != nil {
		temp := o.ExtraData(resp)
		if !judgeInterfaceIsNil(temp) {
			extraData = temp
		}
	}

	if len(delIDs) == 0 && len(changes) == 0 && len(insert) == 0 && !o.Full(resp) && extraData == nil {
		log.ZDebug(o.Ctx, "no data to sync", "table", o.TableName, "entityID", o.EntityID)
		return nil
	}

	if o.Full(resp) {
		err := o.FullSyncer(o.Ctx)
		if err != nil {
			return err
		}
		lvs.UIDList, err = o.FullID(o.Ctx)
		if err != nil {
			return err
		}
	} else {
		if len(delIDs) > 0 {
			lvs.UIDList = datautil.DeleteElems(lvs.UIDList, delIDs...)
		}

		if len(insert) > 0 {
			changes = append(changes, insert...)
		}

		if len(changes) > 0 {
			changeKeys := datautil.SliceSub(datautil.Slice(changes, o.Key), lvs.UIDList)
			if len(changeKeys) > 0 {
				lvs.UIDList = append(lvs.UIDList, changeKeys...)
			}
		}

		local, err := o.Local()
		if err != nil {
			return err
		}

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

		server := datautil.Values(kv)
		if err := o.Syncer(server, local); err != nil {
			return err
		}
		if extraData != nil && o.ExtraDataProcessor != nil {
			if err := o.ExtraDataProcessor(o.Ctx, extraData); err != nil {
				return err
			}
		}

		// The ordering of fullID has changed due to modifications such as group role level changes or friend list reordering.
		// Therefore, it is necessary to refresh and obtain the fullID again.
		if o.IDOrderChanged != nil && o.IDOrderChanged(resp) {
			lvs.UIDList, err = o.FullID(o.Ctx)
			if err != nil {
				return err
			}
		}
	}
	return o.updateVersionInfo(lvs, resp)
}

func (o *VersionSynchronizer[V, R]) CheckVersionSync() error {
	lvs, err := o.getVersionInfo()
	if err != nil {
		return err
	}

	var extraData any
	resp := o.ServerVersion()
	delIDs := o.Delete(resp)
	changes := o.Update(resp)
	insert := o.Insert(resp)
	versionID, version := o.Version(resp)
	if o.ExtraData != nil {
		temp := o.ExtraData(resp)
		if !judgeInterfaceIsNil(temp) {
			extraData = temp
		}
	}

	if len(delIDs) == 0 && len(changes) == 0 && len(insert) == 0 && !o.Full(resp) && extraData == nil {
		log.ZWarn(o.Ctx, "exception no data to sync", errs.New("notification no data"), "table", o.TableName, "entityID", o.EntityID)
		return nil
	}

	log.ZDebug(o.Ctx, "check version sync", "table", o.TableName, "entityID", o.EntityID, "versionID", versionID, "localVersionID", lvs.VersionID, "version", version, "localVersion", lvs.Version)
	/// If the version unique ID cannot correspond with the local version,
	// it indicates that the data might have been tampered with or an exception has occurred.
	//Trigger the complete client-server incremental synchronization.
	if versionID != lvs.VersionID {
		log.ZDebug(o.Ctx, "version id not match", errs.New("version id not match"), "versionID", versionID, "localVersionID", lvs.VersionID)
		o.ServerVersion = nil
		return o.IncrementalSync()
	}

	if lvs.Version+1 == version {
		if len(delIDs) > 0 {
			lvs.UIDList = datautil.DeleteElems(lvs.UIDList, delIDs...)
		}
		if len(insert) > 0 {
			lvs.UIDList = append(lvs.UIDList, datautil.Slice(insert, o.Key)...)
			changes = append(changes, insert...)
		}

		local, err := o.Local()
		if err != nil {
			return err
		}

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

		server := datautil.Values(kv)
		if err := o.Syncer(server, local); err != nil {
			return err
		}
		if extraData != nil && o.ExtraDataProcessor != nil {
			if err := o.ExtraDataProcessor(o.Ctx, extraData); err != nil {
				return err
			}

		}

		// The ordering of fullID has changed due to modifications such as group role level changes or friend list reordering.
		// Therefore, it is necessary to refresh and obtain the fullID again.
		if o.IDOrderChanged != nil && o.IDOrderChanged(resp) {
			lvs.UIDList, err = o.FullID(o.Ctx)
			if err != nil {
				return err
			}
		}
		return o.updateVersionInfo(lvs, resp)
	} else if version <= lvs.Version {
		log.ZWarn(o.Ctx, "version less than local version", errs.New("version less than local version"),
			"table", o.TableName, "entityID", o.EntityID, "version", version, "localVersion", lvs.Version)
		return nil
	} else {
		// If the version number has a gap with the local version number,
		//it indicates that some pushed data might be missing.
		//Trigger the complete client-server incremental synchronization.
		o.ServerVersion = nil
		return o.IncrementalSync()
	}
}
