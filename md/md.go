package md

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"

	"github.com/ggoop/goutils/repositories"
)

const STATE_TEMP = "temp"
const STATE_CREATED = "created"
const STATE_UPDATED = "updated"
const STATE_DELETED = "deleted"
const STATE_NORMAL = "normal"

const TYPE_ENTITY = "entity"
const TYPE_ENUM = "enum"

type MD interface {
	MD() *Mder
}
type Mder struct {
	ID     string
	Name   string
	Domain string
}

type md struct {
	Value interface{}
	db    *repositories.MysqlRepo
}

var initMD bool
var mdCache map[string]*MDEntity

func GetEntity(id string) *MDEntity {
	defer func() {
		if err := recover(); err != nil {
			glog.Error(err)
		}
	}()
	if mdCache == nil {
		mdCache = make(map[string]*MDEntity)
	}
	id = strings.ToLower(id)
	if v, ok := mdCache[id]; ok {
		return v
	}
	item := &MDEntity{}
	if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
		db.Preload("Fields").Take(item, "id=? or code=? or table_name=?", id, id, id)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if item.ID != "" {
		mdCache[strings.ToLower(item.ID)] = item
		mdCache[strings.ToLower(item.Code)] = item
		if item.TableName != "" {
			mdCache[strings.ToLower(item.TableName)] = item
		}
		return item
	}
	return nil
}

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
	if err := m.db.Model(item).Preload("Fields").Take(&item, "id=?", mdInfo.ID).Error; err == nil {
		return &item
	}
	return nil
}
func (m *md) Migrate() {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return
	}
	scope := m.db.NewScope(m.Value)

	entity := m.GetEntity()
	vt := reflect.ValueOf(m.Value).Elem().Type()
	newEntity := &MDEntity{TableName: scope.TableName(), Code: vt.Name(), Name: mdInfo.Name, Domain: mdInfo.Domain}
	newEntity.FullName = vt.String()
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
		if entity.TableName != newEntity.TableName {
			updates["TableName"] = newEntity.TableName
		}
		if entity.Code != newEntity.Code {
			updates["Code"] = newEntity.Code
		}
		if entity.FullName != newEntity.FullName {
			updates["FullName"] = newEntity.FullName
		}
		if entity.Domain != newEntity.Domain {
			updates["Domain"] = newEntity.Domain
		}
		if len(updates) > 0 {
			m.db.Model(entity).Updates(updates)
			entity = m.GetEntity()
		}
	}
	codes := make([]string, 0)
	for _, field := range scope.GetModelStruct().StructFields {
		newField := MDField{Code: field.Name, DBName: field.DBName, IsPrimaryKey: field.IsPrimaryKey, IsNormal: field.IsNormal, Name: field.TagSettings["NAME"], EntityID: entity.ID}
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
		if relationship := field.Relationship; relationship != nil {
			newField.Kind = relationship.Kind
			newField.ForeignKey = strings.Join(relationship.ForeignFieldNames, ".")
			newField.AssociationKey = strings.Join(relationship.AssociationForeignFieldNames, ".")

			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					newField.TypeID = eMd.ID
				}
			}
		}
		newField.Limit = field.TagSettings["LIMIT"]
		if newField.TypeID == "01e9125fe9611c4dc8d47427ea1d5200" || reflectType.Name() == "MDEnum" {
			newField.TypeType = TYPE_ENUM
		} else if newField.TypeID != "" {
			newField.TypeType = TYPE_ENTITY
		} else {
			newField.TypeType = reflectType.Name()
		}
		codes = append(codes, newField.Code)
		oldField := entity.GetField(newField.Code)

		if oldField == nil {
			//新增加
			newField.ID = utils.GUID()
			m.db.Create(&newField)
		} else {
			updates := make(map[string]interface{})
			if oldField.Name != newField.Name {
				updates["Name"] = newField.Name
			}
			if oldField.DBName != newField.DBName {
				updates["DBName"] = newField.DBName
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
			&MDEntity{}, &MDField{}, &MDEnumType{}, &MDEnum{},
		}
		db.AutoMigrate(mds...)
		for _, v := range mds {
			m := newMd(v, db)
			m.Migrate()
		}
	}
	for _, value := range values {
		m := newMd(value, db)
		m.Migrate()
	}
	//表迁移
	db.AutoMigrate(values...)
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
				repo.Table(fmt.Sprintf("%v as t", entity.TableName)).Where(fmt.Sprintf("%v in (?)", field.DBName), ids).Count(&count)
				if count > 0 {
					item := MDEntity{TypeID: entity.TypeID, Code: entity.Code, Name: entity.Name, FullName: entity.FullName, TableName: entity.TableName}
					item.ID = entity.ID
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
