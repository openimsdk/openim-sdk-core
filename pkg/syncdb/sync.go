package syncdb

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

type Complete struct {
	Key  []string
	Data []any
}

// Data 切片中必须是model结构体,字段包含gorm主键和column信息
type Data struct {
	Change   []any       // 更新的数据
	Delete   []any       // 删除的数据
	Complete []*Complete // 完整的数据
}

// model对应的主键和字段
type ModelKey struct {
	PrimaryKey   map[int]string // go model field index -> db column name
	UpdateColumn map[int]string // go model field index -> db column name
}

func NewSync(db *gorm.DB) *Sync {
	return &Sync{
		db:         db,
		modelCache: make(map[string]ModelKey),
	}
}

type Sync struct {
	db         *gorm.DB
	change     []any
	delete     []any
	complete   []Complete
	modelCache map[string]ModelKey
}

func (s *Sync) AddChange(data any) {
	s.change = append(s.change, data)
}

func (s *Sync) AddDelete(data any) {
	s.delete = append(s.delete, data)
}

func (s *Sync) AddComplete(key []string, data []any) {
	s.complete = append(s.complete, Complete{Key: key, Data: data})
}

func (s *Sync) equal(a, b reflect.Value) bool {
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

func (s *Sync) where(m any, primaryKeys map[int]string) *gorm.DB {
	if len(primaryKeys) == 0 {
		panic("no primary key")
	}
	valueOf := reflect.ValueOf(m)
	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	whereDb := s.db
	for index, column := range primaryKeys {
		whereDb = whereDb.Where(fmt.Sprintf("`%s` = ?", column), valueOf.Field(index).Interface())
	}
	return whereDb
}

func (s *Sync) getModelInfo(m any) (*ModelKey, error) {
	if s.modelCache == nil {
		s.modelCache = make(map[string]ModelKey)
	}
	valueOf := reflect.ValueOf(m)
	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return nil, errors.New("not struct slice")
	}
	typeOf := valueOf.Type()
	typeStr := typeOf.String()
	if res, ok := s.modelCache[typeStr]; ok {
		return &res, nil
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
	res := ModelKey{PrimaryKey: primaryKey, UpdateColumn: updateColumn}
	s.modelCache[typeStr] = res
	return &res, nil
}

func (s *Sync) Change() ([]SyncState, error) {
	state := make([]SyncState, 0, len(s.change))
	for i := range s.change {
		change := s.change[i]
		keyInfo, err := s.getModelInfo(change)
		if err != nil {
			return nil, err
		}
		valueOf := reflect.ValueOf(change)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		model := reflect.New(valueOf.Type()) // type: *struct
		if err := s.where(change, keyInfo.PrimaryKey).Take(model.Interface()).Error; err == nil {
			dbData := model.Elem() // type: struct
			changeData := reflect.ValueOf(change)
			for changeData.Kind() == reflect.Pointer {
				changeData = changeData.Elem()
			}
			update := make(map[string]any)
			for index, column := range keyInfo.UpdateColumn {
				changeField := changeData.Field(index)
				if s.equal(dbData.Field(index), changeField) {
					continue
				}
				update[column] = changeField.Interface()
			}
			if len(update) == 0 {
				state = append(state, StateNoChange)
				continue
			}
			if err := s.where(change, keyInfo.PrimaryKey).Model(model.Interface()).Updates(update).Error; err != nil {
				return nil, err
			}
			state = append(state, StateUpdate)
		} else if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(change).Error; err != nil {
				return nil, err
			}
			state = append(state, StateInsert)
		} else {
			return nil, err
		}
	}
	return state, nil
}

func (s *Sync) Delete() ([]SyncState, error) {
	state := make([]SyncState, 0, len(s.delete))
	for i := range s.delete {
		del := s.delete[i]
		keyInfo, err := s.getModelInfo(del)
		if err != nil {
			return nil, err
		}
		valueOf := reflect.ValueOf(del)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		zero := reflect.Zero(valueOf.Type()) // type: struct
		r := s.where(del, keyInfo.PrimaryKey).Delete(zero.Interface())
		if r.Error != nil {
			return nil, r.Error
		}
		if r.RowsAffected == 0 {
			state = append(state, StateNoChange)
		} else {
			state = append(state, StateDelete)
		}
	}
	return state, nil
}

func (s *Sync) checkColumn(cs []any, colIndex []int) error {
	var (
		typeStr string
		pk      = make([]reflect.Value, 0, len(colIndex))
	)
	getColumnValue := func(valueOf reflect.Value) []reflect.Value {
		col := make([]reflect.Value, 0, len(colIndex))
		for _, index := range colIndex {
			col = append(col, valueOf.Field(index))
		}
		return col
	}
	for i := range cs {
		valueOf := reflect.ValueOf(cs[i])
		if valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		if valueOf.Kind() != reflect.Struct {
			return errors.New("not struct")
		}
		if i == 0 {
			pk = getColumnValue(valueOf)
			typeStr = valueOf.Type().String()
			continue
		}
		if valueOf.Type().String() != typeStr {
			return errors.New("not same type")
		}
		for i, value := range getColumnValue(valueOf) {
			if !s.equal(value, pk[i]) {
				return errors.New("not same key")
			}
		}
	}
	return nil
}

func (s *Sync) mapKey(m map[int]string) []int {
	ks := make([]int, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func (s *Sync) checkKey(key []string, primaryKey map[int]string) (map[int]string, error) {
	if len(key)+1 != len(primaryKey) {
		return nil, errors.New("key error")
	}
	col := make(map[int]string)
	if len(key) > 0 {
		temp := make(map[string]struct{})
		for _, k := range key {
			temp[k] = struct{}{}
		}
		if len(key) != len(temp) {
			return nil, errors.New("key error")
		}
		for index, key := range primaryKey {
			if _, ok := temp[key]; ok {
				delete(temp, key)
				col[index] = key
			}
		}
		if len(temp) != 0 {
			return nil, errors.New("key not found")
		}
	}
	return col, nil
}

func (s *Sync) getRecID(m any, indexs []int) string {
	valueOf := reflect.ValueOf(m)
	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	arr := make([]any, 0, len(indexs))
	for _, index := range indexs {
		arr = append(arr, valueOf.Field(index).Interface())
	}
	data, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	t := md5.Sum(data)
	return hex.EncodeToString(t[:])
}

func (s *Sync) Complete() error {
	for i := range s.complete {
		if err := s.completeBy(&s.complete[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sync) completeBy(c *Complete) error {
	if len(c.Data) == 0 {
		return nil
	}
	info, err := s.getModelInfo(c.Data[0])
	if err != nil {
		return err
	}
	// Column
	whereColumn, err := s.checkKey(c.Key, info.PrimaryKey)
	if err != nil {
		return err
	}
	if err := s.checkColumn(c.Data, s.mapKey(whereColumn)); err != nil {
		return err
	}
	first := reflect.ValueOf(c.Data[0])
	for first.Kind() == reflect.Pointer {
		first = first.Elem()
	}
	elemType := first.Type()
	indexs := s.mapKey(info.PrimaryKey)
	idIndex := make(map[string]int)
	for i := range c.Data {
		val := c.Data[i]
		idIndex[s.getRecID(val, indexs)] = i
	}
	if len(idIndex) != len(c.Data) {
		return errors.New("duplicate primary key")
	}
	sliceOf := reflect.SliceOf(elemType)
	const size = 50
	for page := 0; ; page++ {
		dbList := reflect.New(sliceOf)
		dbList.Elem().Set(reflect.MakeSlice(sliceOf, 0, size))
		if err := s.where(first.Interface(), whereColumn).Limit(size).Offset(page * size).Find(dbList.Interface()).Error; err != nil {
			return err
		}
		dbLen := dbList.Elem().Len()
		for i := 0; i < dbLen; i++ {
			item := dbList.Elem().Index(i)
			id := s.getRecID(item.Interface(), indexs)
			idx, ok := idIndex[id]
			if !ok {
				if err := s.where(item.Interface(), info.PrimaryKey).Delete(reflect.New(elemType).Interface()).Error; err != nil {
					return err
				}
				continue
			}
			changeData := reflect.ValueOf(c.Data[idx])
			for changeData.Kind() == reflect.Pointer {
				changeData = changeData.Elem()
			}
			update := make(map[string]any)
			for index, column := range info.UpdateColumn {
				changeField := changeData.Field(index)
				if s.equal(item.Field(index), changeField) {
					continue
				}
				update[column] = changeField.Interface()
			}
			if len(update) == 0 {
				continue
			}
			if err := s.where(changeData.Interface(), info.PrimaryKey).Model(reflect.New(elemType).Interface()).Updates(update).Error; err != nil {
				return err
			}
			delete(idIndex, id)
		}
		if dbLen < size {
			break
		}
	}
	if len(idIndex) > 0 {
		list := reflect.MakeSlice(sliceOf, 0, size)
		for _, i := range idIndex {
			item := c.Data[i]
			val := reflect.ValueOf(item)
			for val.Kind() == reflect.Pointer {
				val = val.Elem()
			}
			list = reflect.Append(list, val)
			if list.Len() >= size {
				if err := s.db.Create(list.Interface()).Error; err != nil {
					return err
				}
				list = reflect.MakeSlice(sliceOf, 0, size)
			}
		}
		if list.Len() > 0 {
			if err := s.db.Create(list.Interface()).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sync) Data() error {

	return nil
}
