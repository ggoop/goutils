package md

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ggoop/goutils/gorm"

	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/repositories"
)

const (
	//临时
	STATE_TEMP = "temp"
	//创建的
	STATE_CREATED = "created"
	//更新的
	STATE_UPDATED = "updated"
	//删除的
	STATE_DELETED = "deleted"
	//正常的
	STATE_NORMAL = "normal"
)
const (
	//简单类型
	TYPE_SIMPLE = "simple"
	//实体
	TYPE_ENTITY = "entity"
	// 枚举
	TYPE_ENUM = "enum"
	// 接口
	TYPE_INTERFACE = "interface"
	// 对象
	TYPE_DTO = "dto"
	// 视图
	TYPE_VIEW = "view"
)

type MD interface {
	MD() *Mder
}
type Mder struct {
	ID     string
	Type   string
	Name   string
	Domain string
}

type md struct {
	Value interface{}
	db    *repositories.MysqlRepo
}

var initMD bool

func newMd(value interface{}, db *repositories.MysqlRepo) *md {
	item := md{Value: value, db: db}
	return &item
}
func (m *md) GetMder() *Mder {
	if mder, ok := m.Value.(MD); ok {
		return mder.MD()
	}
	return nil
}
func (m *md) GetEntity() *MDEntity {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return nil
	}
	item := MDEntity{}
	query := m.db.Model(item).Preload("Fields").Order("id").Where("id=?", mdInfo.ID)
	if err := query.Take(&item).Error; err != nil {
		glog.Error(err)
	} else {
		return &item
	}
	return nil
}

// Get Data Type for MySQL Dialect
func (s *md) dataTypeOf(field *gorm.StructField) string {
	size := 0
	if num, ok := field.TagSettingsGet("SIZE"); ok {
		size, _ = strconv.Atoi(num)
	} else {
		size = 255
	}
	var (
		reflectType = field.Struct.Type
	)

	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	fieldValue := reflect.Indirect(reflect.New(reflectType))
	sqlType := ""

	if sqlType == "" {
		switch fieldValue.Kind() {
		case reflect.Bool:
			sqlType = "boolean"
		case reflect.Int8:
		case reflect.Uint8:
		case reflect.Int, reflect.Int16, reflect.Int32:
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		case reflect.Int64:
		case reflect.Uint64:
			sqlType = "int"
		case reflect.Float32, reflect.Float64:
			sqlType = "decimal"
		case reflect.String:
			if size > 100 {

			}
			sqlType = "string"
		case reflect.Struct:
			if _, ok := fieldValue.Interface().(time.Time); ok {
				sqlType = "datetime"
			}
			if _, ok := fieldValue.Interface().(Time); ok {
				sqlType = "datetime"
			}
			if _, ok := fieldValue.Interface().(decimal.Decimal); ok {
				sqlType = "decimal"
			}
		default:
			sqlType = "string"
		}
	}
	return sqlType
}

func (m *md) Migrate() {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return
	}
	if mdInfo.ID == "" {
		glog.Error("元数据ID为空", glog.String("Name", mdInfo.Name))
		return
	}
	scope := m.db.NewScope(m.Value)

	entity := m.GetEntity()
	vt := reflect.ValueOf(m.Value).Elem().Type()
	newEntity := &MDEntity{TableName: scope.TableName(), Name: mdInfo.Name, Domain: mdInfo.Domain, Code: vt.Name(), Type: mdInfo.Type}
	if newEntity.Type == "" {
		newEntity.Type = TYPE_ENTITY
	}
	if entity == nil {
		entity = newEntity
		entity.ID = mdInfo.ID
		m.db.Create(entity)
		entity = m.GetEntity()
	} else {
		updates := make(map[string]interface{})
		if entity.Name != newEntity.Name {
			updates["Name"] = newEntity.Name
		}
		if entity.Code != newEntity.Code {
			updates["Code"] = newEntity.Code
		}
		if entity.Type != newEntity.Type {
			updates["Type"] = newEntity.Type
		}
		if entity.TableName != newEntity.TableName {
			updates["TableName"] = newEntity.TableName
		}
		if entity.Domain != newEntity.Domain {
			updates["Domain"] = newEntity.Domain
		}
		if len(updates) > 0 {
			m.db.Model(entity).Where("id=?", entity.ID).Updates(updates)
			entity = m.GetEntity()
		}
	}
	if entity == nil {
		glog.Error("元数据ID为空", glog.String("Name", mdInfo.Name))
		return
	}
	codes := make([]string, 0)
	for _, field := range scope.GetModelStruct().StructFields {
		newField := MDField{Code: field.Name, DbName: field.DBName, IsPrimaryKey: field.IsPrimaryKey, IsNormal: field.IsNormal, Name: field.TagSettings["NAME"], EntityID: entity.ID}
		if field.IsIgnored {
			continue
		}

		if newField.Name == "" {
			newField.Name = newField.Code
		}
		//普通数据库字段
		if field.IsNormal {
		}
		reflectType := field.Struct.Type
		if reflectType.Kind() == reflect.Slice {
			reflectType = field.Struct.Type.Elem()
		}
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}
		newField.Limit = field.TagSettings["LIMIT"]
		if relationship := field.Relationship; relationship != nil {
			newField.Kind = relationship.Kind
			newField.ForeignKey = strings.Join(relationship.ForeignFieldNames, ".")
			newField.AssociationKey = strings.Join(relationship.AssociationForeignFieldNames, ".")

			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					newField.TypeID = eMd.ID
					newField.TypeType = eMd.Type
				}
			}
		} else {
			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					newField.TypeID = eMd.ID
					newField.TypeType = eMd.Type
				}
			} else if e := m.dataTypeOf(field); e != "" {
				newField.TypeID = e
			}
		}
		if newField.TypeID != "" && newField.TypeType == "" {
			if typeEntity := GetEntity(newField.TypeID); typeEntity != nil {
				newField.TypeType = typeEntity.Type
			}
		}
		codes = append(codes, newField.Code)
		oldField := entity.GetField(newField.Code)

		if oldField == nil {
			m.db.Create(&newField)
		} else {
			updates := make(map[string]interface{})
			if oldField.Name != newField.Name {
				updates["Name"] = newField.Name
			}
			if oldField.DbName != newField.DbName {
				updates["DbName"] = newField.DbName
			}
			if oldField.AssociationKey != newField.AssociationKey {
				updates["AssociationKey"] = newField.AssociationKey
			}
			if oldField.ForeignKey != newField.ForeignKey {
				updates["ForeignKey"] = newField.ForeignKey
			}
			if oldField.IsNormal != newField.IsNormal {
				updates["IsNormal"] = newField.IsNormal
			}
			if oldField.IsPrimaryKey != newField.IsPrimaryKey {
				updates["IsPrimaryKey"] = newField.IsPrimaryKey
			}
			if oldField.Kind != newField.Kind {
				updates["Kind"] = newField.Kind
			}
			if oldField.TypeID != newField.TypeID {
				updates["TypeID"] = newField.TypeID
			}
			if oldField.TypeType != newField.TypeType {
				updates["TypeType"] = newField.TypeType
			}
			if oldField.Limit != newField.Limit {
				updates["Limit"] = newField.Limit
			}
			if len(updates) > 0 {
				m.db.Model(oldField).Updates(updates)
			}
		}
	}
	//删除不存在的
	if len(entity.Fields) != len(codes) {
		m.db.Where("entity_id=? and code not in (?)", entity.ID, codes).Delete(MDField{})
	}
}

func Migrate(db *repositories.MysqlRepo, values ...interface{}) {
	//先增加模型表
	if !initMD {
		initMD = true
		mds := []interface{}{
			&MDEntity{}, &MDField{}, &MDTag{}, &MDEnum{},
			&MDActionCommand{}, &MDActionRule{}, &MDActionMaker{},
			&MDPage{}, &MDPageView{}, &MDPageWidget{},
			&MDPart{}, &MDPartProps{},
		}
		needDb := make([]interface{}, 0)
		for _, v := range mds {
			m := newMd(v, db)
			if dd := m.GetMder(); dd == nil || dd.Type == TYPE_ENTITY || dd.Type == TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
		}
		db.AutoMigrate(needDb...)
		glog.Error("AutoMigrate MD")
		for _, v := range mds {
			m := newMd(v, db)
			m.Migrate()
		}

		initData(db)
	}
	if len(values) > 0 {
		needDb := make([]interface{}, 0)
		for _, v := range values {
			m := newMd(v, db)
			if dd := m.GetMder(); dd == nil || dd.Type == TYPE_ENTITY || dd.Type == TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
			m.Migrate()
		}
		//表迁移
		glog.Error("AutoMigrate BIZ")
		db.AutoMigrate(needDb...)
	}

}
func QuotedBy(m MD, ids []string, excludes ...MD) ([]MDEntity, []string) {
	if m == nil || ids == nil || len(ids) == 0 {
		return nil, nil
	}
	var repo *repositories.MysqlRepo
	if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
		repo = db
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
		return nil, nil
	}

	excludeIds := make([]string, 0)
	if excludes != nil && len(excludes) > 0 {
		for _, e := range excludes {
			excludeIds = append(excludeIds, e.MD().ID)
		}
	}

	items := make([]MDField, 0)
	query := repo.Table(fmt.Sprintf("%v as f", repo.NewScope(MDField{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v as e on e.id=f.entity_id", repo.NewScope(MDEntity{}).TableName()))
	query = query.Select("f.*")
	if len(excludeIds) > 0 {
		query = query.Where("f.entity_id not in (?)", excludeIds)
	}
	query.Where("f.type_id=? and f.type_type=? and f.kind=?", m.MD().ID, "entity", "belongs_to").Find(&items)
	if len(items) > 0 {
		rtns := make([]MDEntity, 0)
		count := 0
		for _, d := range items {
			entity := GetEntity(d.EntityID)
			if entity == nil || entity.TableName == "" {
				continue
			}
			if d.Kind == "belongs_to" {
				field := entity.GetField(d.ForeignKey)
				if field == nil {
					continue
				}
				repo.Table(fmt.Sprintf("%v as t", entity.TableName)).Where(fmt.Sprintf("%v in (?)", field.DbName), ids).Count(&count)
				if count > 0 {
					item := MDEntity{ID: entity.ID, Type: entity.Type, Name: entity.Name, TableName: entity.TableName}
					rtns = append(rtns, item)
				}
			}
		}
		if len(rtns) > 0 {
			s := make([]string, 0)
			for _, item := range rtns {
				s = append(s, item.Name)
			}
			return rtns, s
		}
	}
	return nil, nil
}
