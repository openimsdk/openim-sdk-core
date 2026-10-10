package datafetcher

import (
	"context"
	"sort"

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
type FetchDataFunc[T any] func(ctx context.Context, uids []string) ([]T, bool, error)

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
	if len(uids) == 0 {
		return nil, nil
	}

	localData, needServer, err := ds.FetchFromLocal(ctx, uids)
	if err != nil {
		return nil, err
	}
	if !needServer {
		return ds.sortByUserIDs(localData, uids), nil
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

	return ds.sortByUserIDs(localData, uids), nil
}

func (ds *DataFetcher[T]) sortByUserIDs(data []T, userIDs []string) []T {
	userIndexMap := make(map[string]int, len(userIDs))
	for i, uid := range userIDs {
		userIndexMap[uid] = i
	}
	sort.SliceStable(data, func(i, j int) bool {
		uid1 := ds.Key(data[i])
		uid2 := ds.Key(data[j])
		index1 := userIndexMap[uid1]
		index2 := userIndexMap[uid2]
		return index1 < index2
	})

	return data
}

// FetchMissingAndCombineLocal fetches missing data from the server and combines it with local data without inserting it into the local database
func (ds *DataFetcher[T]) FetchMissingAndCombineLocal(ctx context.Context, uids []string) ([]T, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	localData, needServer, err := ds.FetchFromLocal(ctx, uids)
	if err != nil {
		return nil, err
	}

	if !needServer {
		return localData, nil
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
		// Combine local data with server data
		localData = append(localData, serverData...)
	}

	return localData, nil
}

func (ds *DataFetcher[T]) FetchWithPaginationV2(ctx context.Context, offset, limit int) ([]T, bool, error) {
	var isEnd bool
	versionInfo, err := ds.db.GetVersionSync(ctx, ds.TableName, ds.EntityID)
	if err != nil {
		return nil, isEnd, err
	}

	if offset > len(versionInfo.UIDList) {
		return nil, isEnd, errs.New("offset exceeds the length of the UID list").Wrap()
	}

	end := offset + limit
	if end >= len(versionInfo.UIDList) {
		isEnd = true
		end = len(versionInfo.UIDList)
	}

	paginatedUIDs := versionInfo.UIDList[offset:end]

	localData, isEnd, err := ds.FetchMissingAndFillLocalV2(ctx, paginatedUIDs, isEnd)
	if err != nil {
		return nil, isEnd, err
	}
	return localData, isEnd, nil
}

func (ds *DataFetcher[T]) FetchMissingAndFillLocalV2(ctx context.Context, uids []string, isEnd bool) ([]T, bool, error) {
	localData, _, err := ds.FetchFromLocal(ctx, uids)
	if err != nil {
		return nil, false, err
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
			return localData, false, nil
		}

		if err := ds.batchInsert(ctx, serverData); err != nil {
			return nil, false, err
		}

		localData = append(localData, serverData...)
	}

	return localData, isEnd, nil
}
