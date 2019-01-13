package md

import (
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

type MD interface {
	MD() *Mder
}
type Mder struct {
	ID   string
	Name string
}

type md struct {
	Value interface{}
	db    *repositories.MysqlRepo
}

var initMD bool
var mdCache map[string]*MDEntity

func GetEntity(id string) *MDEntity {
	if mdCache == nil {
		mdCache = make(map[string]*MDEntity)
	}
	id = strings.ToLower(id)
	if v, ok := mdCache[id]; ok {
		return v
	}
	item := &MDEntity{}
	if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
		db.Preload("Fields").Take(item, "id=? or code=?", id, id)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	if item.ID != "" {
		mdCache[strings.ToLower(item.ID)] = item
		mdCache[strings.ToLower(item.Code)] = item
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
	newEntity := &MDEntity{TableName: scope.TableName(), Code: vt.Name(), Name: mdInfo.Name}
	newEntity.FullName = vt.String()
	if entity == nil {
		entity = newEntity
		entity.ID = mdInfo.ID
		m.db.Create(entity)
		entity = m.GetEntity()
	} else if entity.Name != newEntity.Name || entity.TableName != newEntity.TableName || entity.Code != newEntity.Code || entity.FullName != newEntity.FullName {
		m.db.Model(entity).Update(newEntity)
		entity = m.GetEntity()
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
		if newField.TypeID != "" {
			newField.TypeType = "entity"
		} else {
			newField.TypeType = reflectType.Name()
		}
		codes = append(codes, newField.Code)
		oldField := entity.GetField(newField.Code)

		if oldField == nil {
			//新增加
			newField.ID = utils.GUID()
			m.db.Create(&newField)
		} else if oldField.Name != newField.Name || oldField.DBName != newField.DBName || oldField.AssociationKey != newField.AssociationKey || oldField.ForeignKey != newField.ForeignKey ||
			oldField.IsNormal != newField.IsNormal || oldField.IsPrimaryKey != newField.IsPrimaryKey ||
			oldField.Kind != newField.Kind || oldField.TypeID != newField.TypeID || oldField.TypeType != newField.TypeType {
			//变更的
			m.db.Model(oldField).Update(newField)
		}
	}
	//删除不存在的
	if len(entity.Fields) != len(codes) {
		m.db.Where("entity_id=? and code not in (?)", entity.ID, codes).Delete(MDField{})
	}
}

func Migrate(db *repositories.MysqlRepo, values ...interface{}) {
	if !initMD {
		initMD = true
		//先增加模型表
		db.AutoMigrate(&MDEntity{}, &MDField{}, &MDEnumType{}, &MDEnum{})
		m := newMd(&MDEntity{}, db)
		m.Migrate()

		m = newMd(&MDField{}, db)
		m.Migrate()
	}
	for _, value := range values {
		m := newMd(value, db)
		m.Migrate()
	}
	//表迁移
	db.AutoMigrate(values...)
}
