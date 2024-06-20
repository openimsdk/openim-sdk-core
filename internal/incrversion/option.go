package incrversion

import (
	"context"
	"reflect"
	"sort"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"gorm.io/gorm"
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
}

func (o *VersionSynchronizer[V, R]) getVersionInfo() (*model_struct.LocalVersionSync, error) {
	versionInfo, err := o.DB.GetVersionSync(o.Ctx, o.TableName, o.EntityID)
	if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
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

func (o *VersionSynchronizer[V, R]) Sync() error {
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
			lvs.UIDList = DeleteElements(lvs.UIDList, delIDs)
		}
		if len(insert) > 0 {
			lvs.UIDList = append(lvs.UIDList, datautil.Slice(insert, o.Key)...)

		}
		local, err := o.Local()
		if err != nil {
			return err
		}
		kv := datautil.SliceToMapAny(local, func(v V) (string, V) {
			return o.Key(v), v
		})

		changes = append(changes, insert...)

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
		return o.Sync()
	}
	if lvs.Version+1 == version {
		if len(delIDs) > 0 {
			lvs.UIDList = DeleteElements(lvs.UIDList, delIDs)
		}
		if len(insert) > 0 {
			lvs.UIDList = append(lvs.UIDList, datautil.Slice(insert, o.Key)...)

		}
		local, err := o.Local()
		if err != nil {
			return err
		}
		kv := datautil.SliceToMapAny(local, func(v V) (string, V) {
			return o.Key(v), v
		})
		for i, change := range append(changes, insert...) {
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
		return o.Sync()
	}
}

// DeleteElements 删除切片中包含在另一个切片中的元素，并保持切片顺序
func DeleteElements[E comparable](es []E, toDelete []E) []E {
	// 将要删除的元素存储在哈希集合中
	deleteSet := make(map[E]struct{}, len(toDelete))
	for _, e := range toDelete {
		deleteSet[e] = struct{}{}
	}

	// 通过一个索引 j 来跟踪新的切片位置
	j := 0
	for _, e := range es {
		if _, found := deleteSet[e]; !found {
			es[j] = e
			j++
		}
	}
	return es[:j]
}

// DeleteElement 删除切片中的指定元素，并保持切片顺序
func DeleteElement[E comparable](es []E, element E) []E {
	j := 0
	for _, e := range es {
		if e != element {
			es[j] = e
			j++
		}
	}
	return es[:j]
}

// Slice Converts slice types in batches and sorts the resulting slice using a custom comparator
func Slice[E any, T any](es []E, fn func(e E) T, less func(a, b T) bool) []T {
	// 转换切片
	v := make([]T, len(es))
	for i := 0; i < len(es); i++ {
		v[i] = fn(es[i])
	}

	// 排序切片
	sort.Slice(v, func(i, j int) bool {
		return less(v[i], v[j])
	})

	return v
}
