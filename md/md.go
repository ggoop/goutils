package md

import (
	"reflect"
	"strings"

	"github.com/ggoop/goutils/repositories"
)

type MD interface {
	MD() *MDInfo
}
type MDInfo struct {
	ID   string
	Name string
}

type md struct {
	Value interface{}
	db    *repositories.MysqlRepo
}

func (m *md) GetMDInfo() *MDInfo {
	if mder, ok := m.Value.(MD); ok {
		return mder.MD()
	}
	return nil
}
func (m *md) GetEntity() *Entity {
	mdInfo := m.GetMDInfo()
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
	mdInfo := m.GetMDInfo()
	if mdInfo == nil {
		return
	}
	scope := m.db.NewScope(m.Value)
	entity := m.GetEntity()
	if entity == nil {
		entity = &Entity{TableName: scope.TableName(), Name: mdInfo.Name}
		entity.ID = mdInfo.ID
	}
	for _, field := range scope.GetModelStruct().StructFields {
		entityField := EntityField{Code: field.Name, DBName: field.DBName, IsPrimaryKey: field.IsPrimaryKey, IsNormal: field.IsNormal, Name: field.TagSettings["name"], EntityID: entity.ID}
		if field.IsIgnored {
			continue
		}
		if entityField.Name == "" {
			entityField.Name = entityField.Code
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
			entityField.Kind = relationship.Kind
			entityField.ForeignKey = strings.Join(relationship.ForeignFieldNames, ".")
			entityField.AssociationKey = strings.Join(relationship.AssociationForeignFieldNames, ".")

			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					entityField.TypeID = eMd.ID
				}
			}
		}
	}
}
func Migrate(db *repositories.MysqlRepo, values ...interface{}) {
	//先增加模型表
	db.AutoMigrate(&Entity{}, &EntityField{})
	for _, value := range values {
		m := md{Value: value, db: db}
		m.Migrate()
	}
	//表迁移
	db.AutoMigrate(values...)

}
