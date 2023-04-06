package syncdb

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type SyncState uint8

func (s SyncState) String() string {
	switch s {
	case StateNoChange:
		return "NoChange"
	case StateInsert:
		return "Insert"
	case StateUpdate:
		return "Update"
	case StateDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

const (
	StateNoChange SyncState = 0
	StateInsert   SyncState = 1
	StateUpdate   SyncState = 2
	StateDelete   SyncState = 3
)

type Result struct {
	Change []SyncState
	Delete []SyncState
}

func (r Result) String() string {
	return fmt.Sprintf("Change: %s, Delete: %s", r.Change, r.Delete)
}

// Data 切片中必须是model结构体,字段包含gorm主键和column信息
type Data struct {
	Changes []any // 更新的数据
	Deletes []any // 删除的数据
}

func SyncDB(db *gorm.DB, data *Data) (*Result, error) {
	changeStates, deleteStates, err := syncDB(db, data.Changes, data.Deletes)
	if err != nil {
		return nil, err
	}
	return &Result{
		Change: changeStates,
		Delete: deleteStates,
	}, nil
}

func SyncTxDB(db *gorm.DB, data *Data) (res *Result, err error) {
	err = db.Transaction(func(tx *gorm.DB) (err error) {
		res, err = SyncDB(tx, data)
		return
	})
	return
}

// syncDB 同步数据库 1. 通过主键判断是否存在 2. 通过主键判断是否需要更新
func syncDB(db *gorm.DB, changes []any, deletes []any) (changeStates []SyncState, deleteStates []SyncState, err error) {
	changeStates = make([]SyncState, len(changes))
	deleteStates = make([]SyncState, len(deletes))
	if len(changes) == 0 && len(deletes) == 0 {
		return
	}
	// model对应的主键和字段
	type ModelKey struct {
		PrimaryKey   map[int]string // go model field index -> db column name
		UpdateColumn map[int]string // go model field index -> db column name
	}
	// model字段缓存信息
	modelCache := make(map[string]*ModelKey)
	// 获取model对应的主键和字段
	getModelInfo := func(m any) (*ModelKey, error) {
		valueOf := reflect.ValueOf(m)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		if valueOf.Kind() != reflect.Struct {
			return nil, errors.New("not struct slice")
		}
		typeOf := valueOf.Type()
		typeStr := typeOf.String()
		if res := modelCache[typeStr]; res != nil {
			return res, nil
		}
		var (
			primaryKey   = make(map[int]string)
			updateColumn = make(map[int]string)
		)
		for i := 0; i < typeOf.NumField(); i++ {
			modelTypeField := typeOf.Field(i)
			var (
				column string
				key    bool
			)
			for _, s := range strings.Split(modelTypeField.Tag.Get("gorm"), ";") {
				if strings.HasPrefix(s, "column:") {
					column = s[len("column:"):]
					if column == "-" {
						break
					}
				} else if s == "primary_key" {
					key = true
				}
				if column != "" && key {
					break
				}
			}
			if column == "" {
				return nil, errors.New("no column tag")
			} else if column == "-" {
				continue
			}
			if key {
				primaryKey[i] = column
			} else {
				updateColumn[i] = column
			}
		}
		if len(primaryKey) == 0 {
			return nil, errors.New("no primary key")
		}
		if len(updateColumn) == 0 {
			return nil, errors.New("no update column")
		}
		res := &ModelKey{PrimaryKey: primaryKey, UpdateColumn: updateColumn}
		modelCache[typeStr] = res
		return res, nil
	}
	// 比较值是否相等
	equal := func(a, b reflect.Value) bool {
		for a.Kind() == reflect.Pointer {
			if a.IsNil() && b.IsNil() {
				return true
			}
			if a.IsNil() || b.IsNil() {
				return false
			}
			a = a.Elem()
			b = b.Elem()
		}
		return a.Interface() == b.Interface()
	}
	where := func(m any, primaryKeys map[int]string) *gorm.DB {
		if len(primaryKeys) == 0 {
			panic("no primary key")
		}
		valueOf := reflect.ValueOf(m)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		whereDb := db
		for index, column := range primaryKeys {
			whereDb = whereDb.Where(fmt.Sprintf("`%s` = ?", column), valueOf.Field(index).Interface())
		}
		return whereDb
	}
	for i := range changes {
		change := changes[i]
		keyInfo, err := getModelInfo(change)
		if err != nil {
			return nil, nil, err
		}
		valueOf := reflect.ValueOf(change)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		model := reflect.New(valueOf.Type()) // type: *struct
		if err := where(change, keyInfo.PrimaryKey).Take(model.Interface()).Error; err == nil {
			dbData := model.Elem() // type: struct
			changeData := reflect.ValueOf(change)
			for changeData.Kind() == reflect.Pointer {
				changeData = changeData.Elem()
			}
			update := make(map[string]any)
			for index, column := range keyInfo.UpdateColumn {
				changeField := changeData.Field(index)
				if equal(dbData.Field(index), changeField) {
					changeStates[i] = StateNoChange
					continue
				}
				update[column] = changeField.Interface()
			}
			if err := where(change, keyInfo.PrimaryKey).Model(model.Interface()).Updates(update).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = StateUpdate
		} else if err == gorm.ErrRecordNotFound {
			if err := db.Create(change).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = StateInsert
		} else {
			return nil, nil, err
		}
	}
	for i, model := range deletes {
		del := deletes[i]
		keyInfo, err := getModelInfo(del)
		if err != nil {
			return nil, nil, err
		}
		valueOf := reflect.ValueOf(del)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		zero := reflect.Zero(valueOf.Type()) // type: struct
		res := where(model, keyInfo.PrimaryKey).Delete(zero.Interface())
		if res.Error != nil {
			return nil, nil, res.Error
		}
		if res.RowsAffected == 0 {
			deleteStates[i] = StateNoChange
		} else {
			deleteStates[i] = StateDelete
		}
	}
	return changeStates, deleteStates, nil
}

//func syncDB[T any](db *gorm.DB, changes []*T, deletes []*T) (changeStates []SyncState, deleteStates []SyncState, err error) {
//	changeStates = make([]SyncState, len(changes))
//	deleteStates = make([]SyncState, len(deletes))
//	if len(changes) == 0 && len(deletes) == 0 {
//		return
//	}
//	var zero T
//	modelType := reflect.TypeOf(&zero)
//	for modelType.Kind() == reflect.Pointer {
//		modelType = modelType.Elem()
//	}
//	if modelType.Kind() != reflect.Struct {
//		return nil, nil, errors.New("not struct slice")
//	}
//	var (
//		primaryKey = make(map[int]string) // go field index -> db column name
//		fieldName  = make(map[int]string) // go field index -> db column name
//	)
//	for i := 0; i < modelType.NumField(); i++ {
//		modelTypeField := modelType.Field(i)
//		var (
//			column string
//			key    bool
//		)
//		for _, s := range strings.Split(modelTypeField.Tag.Get("gorm"), ";") {
//			if strings.HasPrefix(s, "column:") {
//				column = s[len("column:"):]
//			} else if s == "primary_key" {
//				key = true
//			}
//		}
//		if column == "" {
//			return nil, nil, errors.New("no column tag")
//		}
//		if key {
//			primaryKey[i] = column
//		} else {
//			fieldName[i] = column
//		}
//	}
//	if len(primaryKey) == 0 {
//		return nil, nil, errors.New("no primary key")
//	}
//	where := func(model *T) *gorm.DB {
//		value := reflect.ValueOf(model)
//		for value.Kind() == reflect.Pointer {
//			value = value.Elem()
//		}
//		whereDb := db
//		for index, column := range primaryKey {
//			whereDb = whereDb.Where(fmt.Sprintf("`%s` = ?", column), value.Field(index).Interface())
//		}
//		return whereDb
//	}
//	equal := func(a, b reflect.Value) bool {
//		for a.Kind() == reflect.Pointer {
//			if a.IsNil() && b.IsNil() {
//				return true
//			}
//			if a.IsNil() || b.IsNil() {
//				return false
//			}
//			a = a.Elem()
//			b = b.Elem()
//		}
//		return a.Interface() == b.Interface()
//	}
//	for i, model := range changes {
//		var t T
//		if err := where(model).Take(&t).Error; err == nil {
//			dbValue := reflect.ValueOf(t)
//			for dbValue.Kind() == reflect.Pointer {
//				dbValue = dbValue.Elem()
//			}
//			changeValue := reflect.ValueOf(model)
//			for changeValue.Kind() == reflect.Pointer {
//				changeValue = changeValue.Elem()
//			}
//			update := make(map[string]any)
//			for index, column := range fieldName {
//				dbFieldValue := dbValue.Field(index)
//				changeFieldValue := changeValue.Field(index)
//				if equal(dbFieldValue, changeFieldValue) {
//					continue
//				}
//				update[column] = changeFieldValue.Interface()
//			}
//			if len(update) == 0 {
//				changeStates[i] = StateNoChange
//				continue
//			}
//			if err := where(model).Model(t).Updates(update).Error; err != nil {
//				return nil, nil, err
//			}
//			changeStates[i] = StateUpdate
//		} else if err == gorm.ErrRecordNotFound {
//			if err := where(model).Create(model).Error; err != nil {
//				return nil, nil, err
//			}
//			changeStates[i] = StateInsert
//		} else {
//			return nil, nil, err
//		}
//	}
//	for i, model := range deletes {
//		var t T
//		res := where(model).Delete(&t)
//		if res.Error != nil {
//			return nil, nil, res.Error
//		}
//		if res.RowsAffected == 0 {
//			deleteStates[i] = StateNoChange
//		} else {
//			deleteStates[i] = StateDelete
//		}
//	}
//	return changeStates, deleteStates, nil
//}
//
//func syncTxDB[T any](db *gorm.DB, changes []*T, deletes []*T) (changeStates []SyncState, deleteStates []SyncState, err error) {
//	err = db.Transaction(func(tx *gorm.DB) (err_ error) {
//		changeStates, deleteStates, err_ = syncDB(tx, changes, deletes)
//		return
//	})
//	return
//}
//
//func SyncTxDB[T any](db any, changes []*T, deletes []*T) (changeStates []SyncState, deleteStates []SyncState, err error) {
//	return syncTxDB(db.(*gorm.DB), changes, deletes)
//}
