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

type State uint8

func (s State) String() string {
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
	StateNoChange State = 0
	StateInsert   State = 1
	StateUpdate   State = 2
	StateDelete   State = 3
)

type Method uint8

func (s Method) String() string {
	switch s {
	case MethodChange:
		return "Change"
	case MethodDelete:
		return "Delete"
	case MethodComplete:
		return "complete"
	default:
		return "Unknown"
	}
}

const (
	MethodChange   = 1
	MethodDelete   = 2
	MethodComplete = 3
)

type Column struct {
	Name  string
	Value any
}

type Where struct {
	Columns []*Column
}

type DBInterface interface {
	// Insert 插入数据 ms 是model层的结构体切片
	Insert(ms any) error // []Struct
	// Delete reflect.Type 返回的是最底层的结构体类型
	Delete(m reflect.Type, where []*Where) error
	// Update reflect.Type 返回的是最底层的结构体类型
	Update(m reflect.Type, where *Where, data map[string]any) error
	// FindOffset reflect.Type 返回的是最底层的结构体类型
	FindOffset(m reflect.Type, where *Where, offset int, limit int) (any, error) // []Struct
}

type complete struct {
	Key   []string // 同于的key
	Data  []any    // 同步的数据
	Value any      // 同步的值
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
	changes        []any
	deletes        []any
	completes      []complete
	modelColumnMap map[string]modelColumnInfo
	size           int
	fn             func(method Method, state State, data any)
}

func (s *Syncer) Listen(fn func(method Method, state State, data any)) *Syncer {
	s.fn = fn
	return s
}

func (s *Syncer) AddChange(data []any) *Syncer {
	s.changes = append(s.changes, data...)
	return s
}

func (s *Syncer) AddDelete(data []any) *Syncer {
	s.deletes = append(s.deletes, data...)
	return s
}

func (s *Syncer) AddComplete(key []string, value any, data []any) *Syncer {
	s.completes = append(s.completes, complete{Key: key, Value: value, Data: data})
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

// Change 变更数据
func (s *Syncer) Change() error {
	for i := range s.changes {
		srvItem := s.changes[i]
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
			s.on(MethodChange, StateInsert, srvItem)
		} else { // 存在
			update := s.getModelStructUpdate(dbRes.Index(0), valueOf, keyInfo.UpdateColumn)
			if len(update) > 0 {
				if err := s.db.Update(valueOf.Type(), where, update); err != nil {
					return err
				}
				s.on(MethodChange, StateUpdate, srvItem)
			} else {
				s.on(MethodChange, StateNoChange, srvItem)
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
			s.on(MethodDelete, StateNoChange, del) // 不存在
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

func (s *Syncer) Complete() error {
	if s.size <= 0 {
		return errors.New("size must > 0")
	}
	for i := range s.completes {
		if err := s.completeBy(&s.completes[i]); err != nil {
			return err
		}
	}
	return nil
}

// checkColumn 全量同步, 检查类型是否一致, 指定的列的值是否相等
func (s *Syncer) checkColumn(value reflect.Value, cs []any, columnFieldIndex []int) error {
	getColumnValue := func(valueOf reflect.Value) []reflect.Value {
		column := make([]reflect.Value, 0, len(columnFieldIndex))
		for _, index := range columnFieldIndex {
			column = append(column, valueOf.Field(index))
		}
		return column
	}
	typeStr := value.Type().String()
	firstColumnValue := getColumnValue(value)
	for i := range cs {
		valueOf := reflect.ValueOf(cs[i])
		if valueOf.Kind() == reflect.Pointer {
			valueOf = valueOf.Elem()
		}
		if valueOf.Kind() != reflect.Struct {
			return errors.New("not struct")
		}
		if valueOf.Type().String() != typeStr {
			return errors.New("not same type")
		}
		for i, value := range getColumnValue(valueOf) {
			if !s.equal(value, firstColumnValue[i]) {
				return errors.New("not same value")
			}
		}
	}
	return nil
}

// completeBy 传入完成的数据，根据指定主键列新增,删除,修改数据
func (s *Syncer) completeBy(c *complete) error {
	if c.Value == nil {
		return errors.New("value is nil")
	}
	baseValueOf := reflect.ValueOf(c.Value)
	for baseValueOf.Kind() == reflect.Pointer {
		baseValueOf = baseValueOf.Elem()
	}
	if baseValueOf.Kind() != reflect.Struct {
		return fmt.Errorf("not struct %s", baseValueOf.Kind().String())
	}
	columnInfo, err := s.getModelColumnInfo(baseValueOf.Interface())
	if err != nil {
		return err
	}
	whereColumnInfo, err := s.getWhereKey(c.Key, columnInfo.PrimaryKey)
	if err != nil {
		return err
	}
	if err := s.checkColumn(baseValueOf, c.Data, s.mapKey(whereColumnInfo)); err != nil {
		return err
	}
	elemTypeOf := baseValueOf.Type()
	primaryKeyIndexes := s.mapKey(columnInfo.PrimaryKey) // 主键的索引
	recordIDIndexMap := make(map[string]int)             // 记录主键生成的唯一ID 和 数据在Data中的索引
	for i := range c.Data {
		recordIDIndexMap[s.getRecordID(reflect.ValueOf(c.Data[i]), primaryKeyIndexes)] = i
	}
	if len(recordIDIndexMap) != len(c.Data) {
		return errors.New("gen primary key not unique")
	}
	where := s.getColumnWhere(baseValueOf, whereColumnInfo) // 需要同步的查询条件
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
		n := dbResValueOf.Len() // 查询的数量
		for i := 0; i < n; i++ {
			dbDataValueOf := dbResValueOf.Index(i)
			for dbDataValueOf.Kind() == reflect.Pointer {
				dbDataValueOf = dbDataValueOf.Elem()
			}
			if elemTypeOf.String() != dbDataValueOf.Type().String() {
				return fmt.Errorf("not same type %s %s", elemTypeOf.String(), dbDataValueOf.Type().String()) // 类型不匹配
			}
			id := s.getRecordID(dbDataValueOf, primaryKeyIndexes)
			idx, ok := recordIDIndexMap[id]
			if !ok {
				if err := s.db.Delete(elemTypeOf, []*Where{s.getColumnWhere(dbDataValueOf, columnInfo.PrimaryKey)}); err != nil {
					return err
				}
				continue
			}
			srvData := reflect.ValueOf(c.Data[idx])
			for srvData.Kind() == reflect.Pointer {
				srvData = srvData.Elem()
			}
			update := s.getModelStructUpdate(dbDataValueOf, srvData, columnInfo.UpdateColumn)
			if len(update) == 0 {
				s.on(MethodComplete, StateNoChange, dbDataValueOf.Interface())
				continue
			}
			if err := s.db.Update(elemTypeOf, s.getColumnWhere(dbDataValueOf, columnInfo.PrimaryKey), update); err != nil {
				return err
			}
			delete(recordIDIndexMap, id) // 删除已经更新的数据索引
			s.on(MethodComplete, StateUpdate, dbDataValueOf.Interface())
		}
		if n < s.size {
			break
		}
	}
	if len(recordIDIndexMap) > 0 {
		inserts := reflect.MakeSlice(reflect.SliceOf(elemTypeOf), 0, s.size) // 分批插入需要的切片
		for _, index := range recordIDIndexMap {
			inserts = reflect.Append(inserts, reflect.ValueOf(c.Data[index]))
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
		for _, i := range recordIDIndexMap {
			s.on(MethodComplete, StateInsert, c.Data[i])
		}
	}
	return nil
}

func (s *Syncer) Start() error {
	if err := s.Change(); err != nil {
		return err
	}
	if err := s.Delete(); err != nil {
		return err
	}
	if err := s.Complete(); err != nil {
		return err
	}
	return nil
}

func (s *Syncer) on(method Method, state State, data any) {
	if s.fn != nil {
		s.fn(method, state, data)
	}
}
