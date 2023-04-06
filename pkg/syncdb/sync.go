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

func SyncDB[T any](db *gorm.DB, changes []*T, deletes []*T) (changeStates []SyncState, deleteStates []SyncState, err error) {
	changeStates = make([]SyncState, len(changes))
	deleteStates = make([]SyncState, len(deletes))
	if len(changes) == 0 && len(deletes) == 0 {
		return
	}
	var zero T
	modelType := reflect.TypeOf(&zero)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil, nil, errors.New("not struct slice")
	}
	var (
		primaryKey = make(map[int]string) // go field index -> db column name
		fieldName  = make(map[int]string) // go field index -> db column name
	)
	for i := 0; i < modelType.NumField(); i++ {
		modelTypeField := modelType.Field(i)
		var (
			column string
			key    bool
		)
		for _, s := range strings.Split(modelTypeField.Tag.Get("gorm"), ";") {
			if strings.HasPrefix(s, "column:") {
				column = s[len("column:"):]
			} else if s == "primary_key" {
				key = true
			}
		}
		if column == "" {
			return nil, nil, errors.New("no column tag")
		}
		if key {
			primaryKey[i] = column
		} else {
			fieldName[i] = column
		}
	}
	if len(primaryKey) == 0 {
		return nil, nil, errors.New("no primary key")
	}
	where := func(model *T) *gorm.DB {
		value := reflect.ValueOf(model)
		for value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		whereDb := db
		for index, column := range primaryKey {
			whereDb = whereDb.Where(fmt.Sprintf("`%s` = ?", column), value.Field(index).Interface())
		}
		return whereDb
	}
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
	for i, model := range changes {
		var t T
		if err := where(model).Take(&t).Error; err == nil {
			dbValue := reflect.ValueOf(t)
			for dbValue.Kind() == reflect.Pointer {
				dbValue = dbValue.Elem()
			}
			changeValue := reflect.ValueOf(model)
			for changeValue.Kind() == reflect.Pointer {
				changeValue = changeValue.Elem()
			}
			update := make(map[string]any)
			for index, column := range fieldName {
				dbFieldValue := dbValue.Field(index)
				changeFieldValue := changeValue.Field(index)
				if equal(dbFieldValue, changeFieldValue) {
					continue
				}
				update[column] = changeFieldValue.Interface()
			}
			if len(update) == 0 {
				changeStates[i] = StateNoChange
				continue
			}
			if err := where(model).Model(t).Updates(update).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = StateUpdate
		} else if err == gorm.ErrRecordNotFound {
			if err := where(model).Create(model).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = StateInsert
		} else {
			return nil, nil, err
		}
	}
	for i, model := range deletes {
		var t T
		res := where(model).Delete(&t)
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

func SyncTxDB[T any](db *gorm.DB, changes []*T, deletes []*T) (changeStates []SyncState, deleteStates []SyncState, err error) {
	err = db.Transaction(func(tx *gorm.DB) (err_ error) {
		changeStates, deleteStates, err_ = SyncDB(tx, changes, deletes)
		return
	})
	return
}
