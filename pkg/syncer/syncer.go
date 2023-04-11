package syncer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type globalItem struct {
	Key  map[string]any // 同于的key
	Data any            // 同步的数据
	Safe bool           // 为true是,会检查key中指定列的值与data中的值是否相等,如果不相等,则会返回error
}

// modelColumnInfo model对应的主键和字段
type modelColumnInfo struct {
	PrimaryKey   map[int]string // go model field index -> db column name
	UpdateColumn map[int]string // go model field index -> db column name
}

func New(db DBInterface) *Syncer {
	return &Syncer{
		db:             db,
		size:           50,
		modelColumnMap: make(map[string]modelColumnInfo),
	}
}

type Syncer struct {
	db             DBInterface
	locally        []any
	deletes        []any
	globals        []globalItem // 全局数据
	modelColumnMap map[string]modelColumnInfo
	size           int
	fn             func(method Method, state State, data any)
}

func (s *Syncer) Listen(fn func(method Method, state State, data any)) *Syncer {
	s.fn = fn
	return s
}

func (s *Syncer) AddLocally(data []any) *Syncer {
	s.locally = append(s.locally, data...)
	return s
}

func (s *Syncer) AddDelete(data []any) *Syncer {
	s.deletes = append(s.deletes, data...)
	return s
}

func (s *Syncer) AddGlobal(key map[string]any, data any) *Syncer {
	s.globals = append(s.globals, globalItem{Key: key, Data: data, Safe: true})
	return s
}

// equal 同类型比较是否相等
func (s *Syncer) equal(a, b reflect.Value) bool {
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

// getModelColumnInfo 获取model对应的主键和字段
func (s *Syncer) getModelColumnInfo(m any) (*modelColumnInfo, error) {
	valueOf := reflect.ValueOf(m)
	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return nil, errors.New("not struct slice")
	}
	typeOf := valueOf.Type()
	typeStr := typeOf.String()
	if res, ok := s.modelColumnMap[typeStr]; ok {
		return &res, nil
	}
	var (
		primaryKey   = make(map[int]string) // go model struct field index -> db column name
		updateColumn = make(map[int]string) // go model struct field index -> db column name
	)
	for i := 0; i < typeOf.NumField(); i++ {
		modelTypeField := typeOf.Field(i)
		var (
			column string // 数据库字段名
			key    bool   // 是否主键
		)
		// 解析gorm tag
		for _, s := range strings.Split(modelTypeField.Tag.Get("gorm"), ";") {
			if strings.HasPrefix(s, "column:") {
				column = s[len("column:"):]
				if column == "-" { // gorm 忽略字段
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
	res := modelColumnInfo{PrimaryKey: primaryKey, UpdateColumn: updateColumn}
	s.modelColumnMap[typeStr] = res
	return &res, nil
}

// getColumnWhere 根据列信息获取where条件
func (s *Syncer) getColumnWhere(m reflect.Value, columnInfoMap map[int]string) *Where {
	where := &Where{Columns: make([]*Column, 0, len(columnInfoMap))}
	for index, column := range columnInfoMap {
		where.Columns = append(where.Columns, &Column{Name: column, Value: m.Field(index).Interface()})
	}
	return where
}

// getModelStructUpdate 获取需要更新的字段
func (s *Syncer) getModelStructUpdate(dbStruct reflect.Value, srvStruct reflect.Value, updateColumn map[int]string) map[string]any {
	update := make(map[string]any)
	for index, column := range updateColumn {
		field := srvStruct.Field(index)
		if s.equal(dbStruct.Field(index), field) {
			continue
		}
		update[column] = field.Interface()
	}
	return update
}

// Locally 变更数据
func (s *Syncer) Locally() error {
	for i := range s.locally {
		srvItem := s.locally[i]
		keyInfo, err := s.getModelColumnInfo(srvItem)
		if err != nil {
			return err
		}
		valueOf := reflect.ValueOf(srvItem)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		where := s.getColumnWhere(valueOf, keyInfo.PrimaryKey)
		models, err := s.db.FindOffset(valueOf.Type(), where, 0, 1) // 根据主键查询 返回值是一个切片
		if err != nil {
			return err
		}
		dbRes := reflect.ValueOf(models)
		for dbRes.Kind() == reflect.Pointer {
			dbRes = dbRes.Elem()
		}
		elemType := dbRes.Type().Elem()
		for elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}
		if elemType.String() != valueOf.Type().String() {
			return fmt.Errorf("model type not equal %s != %s", elemType.String(), valueOf.Type().String())
		}
		if dbRes.Len() == 0 { // 不存在
			slice := reflect.MakeSlice(reflect.SliceOf(valueOf.Type()), 1, 1) // 创建一个切片
			slice.Index(0).Set(valueOf)
			if err := s.db.Insert(slice.Interface()); err != nil {
				return err
			}
			s.on(MethodLocally, StateInsert, srvItem)
		} else { // 存在
			update := s.getModelStructUpdate(dbRes.Index(0), valueOf, keyInfo.UpdateColumn)
			if len(update) > 0 {
				if err := s.db.Update(valueOf.Type(), where, update); err != nil {
					return err
				}
				s.on(MethodLocally, StateUpdate, srvItem)
			} else {
				s.on(MethodLocally, StateUnchanged, srvItem)
			}
		}
	}
	return nil
}

// Delete 删除数据
func (s *Syncer) Delete() error {
	for i := range s.deletes {
		del := s.deletes[i]
		keyInfo, err := s.getModelColumnInfo(del)
		if err != nil {
			return err
		}
		valueOf := reflect.ValueOf(del)
		for valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		where := s.getColumnWhere(valueOf, keyInfo.PrimaryKey)
		dbRes, err := s.db.FindOffset(valueOf.Type(), where, 0, 1)
		if err != nil {
			return err
		}
		dbResValueOf := reflect.ValueOf(dbRes)
		for dbResValueOf.Kind() == reflect.Pointer {
			dbResValueOf = dbResValueOf.Elem()
		}
		if dbResValueOf.Kind() != reflect.Slice {
			return fmt.Errorf("model type not slice, %s", dbResValueOf.Kind().String())
		}
		if dbResValueOf.Len() == 0 {
			s.on(MethodDelete, StateUnchanged, del) // 不存在
		} else {
			if err := s.db.Delete(valueOf.Type(), []*Where{where}); err != nil {
				return err
			}
			itemValueOf := dbResValueOf.Index(0)
			for itemValueOf.Kind() == reflect.Pointer {
				itemValueOf = itemValueOf.Elem()
			}
			s.on(MethodDelete, StateDelete, itemValueOf.Interface())
		}
	}
	return nil
}

// mapKey 返回map的key
func (s *Syncer) mapKey(m map[int]string) []int {
	ks := make([]int, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

// getWhereKey 检查key是否存在, 且只能为逐渐key, 比主键数量少1
func (s *Syncer) getWhereKey(key []string, primaryKey map[int]string) (map[int]string, error) {
	if len(key)+1 != len(primaryKey) {
		return nil, errors.New("key error")
	}
	column := make(map[int]string)
	if len(key) > 0 {
		temp := make(map[string]struct{})
		for _, k := range key {
			temp[k] = struct{}{}
		}
		if len(key) != len(temp) {
			return nil, fmt.Errorf("primary key %d curres %d", len(primaryKey), len(temp))
		}
		for index, key := range primaryKey {
			if _, ok := temp[key]; ok {
				delete(temp, key)
				column[index] = key
			}
		}
		if len(temp) != 0 {
			return nil, errors.New("key not found")
		}
	}
	return column, nil
}

// getRecordID 获取记录多个主键生成的唯一ID
func (s *Syncer) getRecordID(valueOf reflect.Value, indexes []int) string {
	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	arr := make([]any, 0, len(indexes))
	for _, index := range indexes {
		arr = append(arr, valueOf.Field(index).Interface())
	}
	data, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	t := md5.Sum(data)
	return hex.EncodeToString(t[:])
}

func (s *Syncer) Global() error {
	if s.size <= 0 {
		return errors.New("size must > 0")
	}
	for i := range s.globals {
		if err := s.globalItem(&s.globals[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Syncer) Start() error {
	if err := s.Locally(); err != nil {
		return err
	}
	if err := s.Delete(); err != nil {
		return err
	}
	if err := s.Global(); err != nil {
		return err
	}
	return nil
}

func (s *Syncer) on(method Method, state State, data any) {
	if s.fn != nil {
		s.fn(method, state, data)
	}
}

func (s *Syncer) getWhere(key map[string]any, info *modelColumnInfo, typeOf reflect.Type) (map[int]string, error) {
	if len(key) >= len(info.PrimaryKey) {
		return nil, errors.New("key len error")
	}
	res := make(map[int]string)
	if len(key) == 0 {
		return res, nil
	}
	columns := make(map[string]int)
	for index, column := range info.PrimaryKey {
		columns[column] = index
	}
	for column, value := range key {
		index, ok := columns[column]
		if !ok {
			return nil, fmt.Errorf("column %s not found", column)
		}
		fieldTypeOf := typeOf.Field(index)
		if reflect.TypeOf(value).String() != fieldTypeOf.Type.String() {
			return nil, fmt.Errorf("column %s type error", column)
		}
		res[index] = column
	}
	return res, nil
}

// globalItem 传入完成的数据，根据指定主键列新增,删除,修改数据
func (s *Syncer) globalItem(c *globalItem) error {
	svrValueOf := reflect.ValueOf(c.Data)
	for svrValueOf.Kind() == reflect.Pointer {
		svrValueOf = svrValueOf.Elem()
	}
	if svrValueOf.Kind() != reflect.Slice {
		return errors.New("not slice")
	}
	elemTypeOf := svrValueOf.Type().Elem()
	for elemTypeOf.Kind() == reflect.Pointer {
		elemTypeOf = elemTypeOf.Elem()
	}
	elemValueOf := reflect.New(elemTypeOf).Elem()
	columnInfo, err := s.getModelColumnInfo(elemValueOf.Interface())
	if err != nil {
		return err
	}
	if c.Key == nil {
		c.Key = make(map[string]any)
	}
	whereColumnInfo, err := s.getWhere(c.Key, columnInfo, elemTypeOf)
	if err != nil {
		return err
	}
	for index, column := range whereColumnInfo {
		elemValueOf.Field(index).Set(reflect.ValueOf(c.Key[column]))
	}
	where := s.getColumnWhere(elemValueOf, whereColumnInfo)
	srvRecordIDIndexMap := make(map[string]int) // 记录主键生成的唯一ID 和 数据在Data中的索引
	primaryKeyIndexes := s.mapKey(columnInfo.PrimaryKey)
	for i := svrValueOf.Len() - 1; i >= 0; i-- {
		itemValueOf := svrValueOf.Index(i)
		if c.Safe && len(c.Key) > 0 {
			for index := range whereColumnInfo {
				if !s.equal(elemValueOf.Field(index), itemValueOf.Field(index)) { // 比较指定的key对应的值是否相等
					return errors.New("safe key not equal")
				}
			}
		}
		srvRecordIDIndexMap[s.getRecordID(itemValueOf, primaryKeyIndexes)] = i
	}
	for page := 0; ; page++ {
		dbRes, err := s.db.FindOffset(elemTypeOf, where, page*s.size, s.size) // 查询数据库中的数据 返回值为切片
		if err != nil {
			return err
		}
		dbResValueOf := reflect.ValueOf(dbRes)
		for dbResValueOf.Kind() == reflect.Pointer {
			dbResValueOf = dbResValueOf.Elem()
		}
		if dbResValueOf.Kind() != reflect.Slice {
			return fmt.Errorf("not slice %s", dbResValueOf.Kind().String())
		}
		itemTypeOf := dbResValueOf.Type().Elem()
		for elemTypeOf.Kind() == reflect.Pointer {
			elemTypeOf = elemTypeOf.Elem()
		}
		if elemTypeOf.String() != itemTypeOf.String() {
			return fmt.Errorf("not same type %s %s", elemTypeOf.String(), itemTypeOf.String()) // 类型不匹配
		}
		n := dbResValueOf.Len() // 查询的数量
		for i := 0; i < n; i++ {
			dbDataValueOf := dbResValueOf.Index(i)
			for dbDataValueOf.Kind() == reflect.Pointer {
				dbDataValueOf = dbDataValueOf.Elem()
			}
			id := s.getRecordID(dbDataValueOf, primaryKeyIndexes)
			idx, ok := srvRecordIDIndexMap[id]
			if !ok {
				if err := s.db.Delete(elemTypeOf, []*Where{s.getColumnWhere(dbDataValueOf, columnInfo.PrimaryKey)}); err != nil {
					return err
				}
				continue
			}
			srvData := svrValueOf.Index(idx)
			for srvData.Kind() == reflect.Pointer {
				srvData = srvData.Elem()
			}
			update := s.getModelStructUpdate(dbDataValueOf, srvData, columnInfo.UpdateColumn)
			if len(update) == 0 {
				s.on(MethodGlobal, StateUnchanged, dbDataValueOf.Interface())
				continue
			}
			if err := s.db.Update(elemTypeOf, s.getColumnWhere(dbDataValueOf, columnInfo.PrimaryKey), update); err != nil {
				return err
			}
			delete(srvRecordIDIndexMap, id) // 删除已经更新的数据索引
			s.on(MethodGlobal, StateUpdate, dbDataValueOf.Interface())
		}
		if n < s.size {
			break
		}
	}
	if len(srvRecordIDIndexMap) > 0 {
		inserts := reflect.MakeSlice(reflect.SliceOf(elemTypeOf), 0, s.size) // 分批插入需要的切片
		for _, index := range srvRecordIDIndexMap {
			inserts = reflect.Append(inserts, svrValueOf.Index(index))
			if inserts.Len() >= s.size {
				if err := s.db.Insert(inserts.Interface()); err != nil {
					return err
				}
				inserts = reflect.MakeSlice(reflect.SliceOf(elemTypeOf), 0, s.size)
			}
		}
		if inserts.Len() > 0 {
			if err := s.db.Insert(inserts.Interface()); err != nil {
				return err
			}
		}
		for _, i := range srvRecordIDIndexMap {
			s.on(MethodGlobal, StateInsert, svrValueOf.Index(i))
		}
	}
	return nil
}
