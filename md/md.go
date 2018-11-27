package md

import (
	"reflect"
	"strings"

	"github.com/ggoop/goutils/repositories"
)

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
func (m *md) GetEntity() *Entity {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return nil
	}
	item := Entity{}
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
	newEntity := &Entity{TableName: scope.TableName(), Name: mdInfo.Name}
	if entity == nil {
		entity = newEntity
		entity.ID = mdInfo.ID
		m.db.Create(entity)
		entity = m.GetEntity()
	} else if entity.Name != newEntity.Name || entity.TableName != newEntity.TableName {
		m.db.Model(entity).Update(newEntity)
		entity = m.GetEntity()
	}
	codes := make([]string, 0)
	for _, field := range scope.GetModelStruct().StructFields {
		newField := EntityField{Code: field.Name, DBName: field.DBName, IsPrimaryKey: field.IsPrimaryKey, IsNormal: field.IsNormal, Name: field.TagSettings["NAME"], EntityID: entity.ID}
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
		codes = append(codes, newField.Code)
		oldField := entity.GetField(newField.Code)

		if oldField == nil {
			//新增加
			m.db.Create(&newField)
		} else if oldField.Name != newField.Name || oldField.DBName != newField.DBName || oldField.AssociationKey != newField.AssociationKey || oldField.ForeignKey != newField.ForeignKey ||
			oldField.IsNormal != newField.IsNormal || oldField.IsPrimaryKey != newField.IsPrimaryKey ||
			oldField.Kind != newField.Kind || oldField.TypeID != newField.TypeID {
			//变更的
			m.db.Model(oldField).Update(newField)
		}
	}
	//删除不存在的
	if len(entity.Fields) != len(codes) {
		m.db.Where("entity_id=? and code not in (?)", entity.ID, codes).Delete(EntityField{})
	}
}
func Migrate(db *repositories.MysqlRepo, values ...interface{}) {
	//先增加模型表
	db.AutoMigrate(&Entity{}, &EntityField{})
	for _, value := range values {
		m := newMd(value, db)
		m.Migrate()
	}
	//表迁移
	db.AutoMigrate(values...)

}
