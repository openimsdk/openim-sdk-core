package db_interface

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func SyncDB[T any](db *gorm.DB, changes []T, deletes []T) ([]uint8, []uint8, error) {
	if len(changes) == 0 && len(deletes) == 0 {
		return []uint8{}, []uint8{}, nil
	}
	var zero T
	modelType := reflect.TypeOf(&zero)
	for modelType.Kind() == reflect.Ptr {
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
	where := func(model T) *gorm.DB {
		value := reflect.ValueOf(model)
		for value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		whereDb := db
		for index, column := range primaryKey {
			whereDb = whereDb.Where(fmt.Sprintf("`%s` = ?", column), value.Field(index).Interface())
		}
		return whereDb
	}
	equal := func(a, b reflect.Value) bool {
		for a.Kind() == reflect.Ptr {
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
	changeStates := make([]uint8, len(changes)) // 0: no change, 1: update, 2: insert
	for i, model := range changes {
		var t T
		if err := where(model).Take(&t).Error; err == nil {
			dbValue := reflect.ValueOf(t)
			for dbValue.Kind() == reflect.Ptr {
				dbValue = dbValue.Elem()
			}
			changeValue := reflect.ValueOf(model)
			for dbValue.Kind() == reflect.Ptr {
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
				changeStates[i] = 0
				continue
			}
			if err := where(model).Model(t).Updates(update).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = 1
		} else if err == gorm.ErrRecordNotFound {
			if err := where(model).Create(model).Error; err != nil {
				return nil, nil, err
			}
			changeStates[i] = 2
		} else {
			return nil, nil, err
		}
	}
	deleteStates := make([]uint8, len(deletes)) // 0: no change, 1: delete
	for i, model := range deletes {
		var t T
		res := where(model).Delete(&t)
		if res.Error != nil {
			return nil, nil, res.Error
		}
		if res.RowsAffected == 0 {
			deleteStates[i] = 0
		} else {
			deleteStates[i] = 1
		}
	}
	return changeStates, deleteStates, nil
}
