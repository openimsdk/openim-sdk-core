package datafetcher

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
)

// DataFetcher is a struct that handles data synchronization
type DataFetcher[T any] struct {
	db              db_interface.VersionSyncModel
	TableName       string
	EntityID        string
	Key             func(T) string
	batchInsert     func(ctx context.Context, servers []T) error
	FetchFromLocal  FetchDataFunc[T]
	fetchFromServer FetchFromServerFunc[T]
}

// FetchDataFunc is a function type for fetching data
type FetchDataFunc[T any] func(ctx context.Context, uids []string) ([]T, error)

// FetchFromServerFunc is a function type for fetching data from server
type FetchFromServerFunc[T any] func(ctx context.Context, uids []string) ([]T, error)

// NewDataFetcher creates a new NewDataFetcher
func NewDataFetcher[T any](db db_interface.VersionSyncModel, tableName string, entityID string, key func(T) string,
	batchInsert func(ctx context.Context, servers []T) error, fetchFromLocal FetchDataFunc[T], fetchFromServer FetchFromServerFunc[T]) *DataFetcher[T] {
	return &DataFetcher[T]{
		db:              db,
		TableName:       tableName,
		EntityID:        entityID,
		Key:             key,
		batchInsert:     batchInsert,
		FetchFromLocal:  fetchFromLocal,
		fetchFromServer: fetchFromServer,
	}
}

// FetchWithPagination fetches data with pagination and fills missing data from server
func (ds *DataFetcher[T]) FetchWithPagination(ctx context.Context, offset, limit int) ([]T, error) {
	versionInfo, err := ds.db.GetVersionSync(ctx, ds.TableName, ds.EntityID)
	if err != nil {
		return nil, err
	}

	if offset > len(versionInfo.UIDList) {
		return nil, errs.New("offset exceeds the length of the UID list").Wrap()
	}

	end := offset + limit
	if end > len(versionInfo.UIDList) {
		end = len(versionInfo.UIDList)
	}

	paginatedUIDs := versionInfo.UIDList[offset:end]

	localData, err := ds.FetchMissingAndFillLocal(ctx, paginatedUIDs)
	if err != nil {
		return nil, err
	}

	return localData, nil
}

// FetchMissingAndFillLocal fetches missing data from server and fills local database
func (ds *DataFetcher[T]) FetchMissingAndFillLocal(ctx context.Context, uids []string) ([]T, error) {
	localData, err := ds.FetchFromLocal(ctx, uids)
	if err != nil {
		return nil, err
	}

	localUIDSet := datautil.SliceSetAny(localData, ds.Key)

	var missingUIDs []string
	for _, uid := range uids {
		if _, found := localUIDSet[uid]; !found {
			missingUIDs = append(missingUIDs, uid)
		}
	}

	if len(missingUIDs) > 0 {
		serverData, err := ds.fetchFromServer(ctx, missingUIDs)
		if err != nil {
			return nil, err
		}
		if len(serverData) > 0 {
			if err := ds.batchInsert(ctx, serverData); err != nil {
				return nil, err
			}

			localData = append(localData, serverData...)
		}

	}

	return localData, nil
}
