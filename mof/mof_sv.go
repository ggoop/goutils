package mof

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

type MOFSv struct {
	repo *repositories.MysqlRepo
}

func NewMOFSv(repo *repositories.MysqlRepo) *MOFSv {
	return &MOFSv{repo: repo}
}

func (s *MOFSv) getEntity(entityID string) md.MDEntity {
	item := md.MDEntity{}
	s.repo.Model(item).Order("id").Take(&item, "id=?", entityID)
	return item
}
func (s *MOFSv) entityToTables(items []md.MDEntity, oldItems []md.MDEntity) error {
	if items == nil || len(items) == 0 {
		return nil
	}
	for i, _ := range items {
		entity := items[i]
		oldItem := md.MDEntity{}
		for oi, ov := range oldItems {
			if entity.ID == ov.ID {
				oldItem = oldItems[oi]
				break
			}
		}
		if entity.Type != md.TYPE_ENTITY {
			continue
		}
		if entity.TableName == "" {
			entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
		}
		if s.repo.Dialect().HasTable(entity.TableName) {
			s.updateTable(entity, oldItem)
		} else {
			s.createTable(entity)
		}
	}
	return nil
}
func (s *MOFSv) createTable(item md.MDEntity) {
	if len(item.Fields) == 0 {
		return
	}
	var primaryKeys []string
	var tags []string
	for i, _ := range item.Fields {
		field := item.Fields[i]
		if field.IsNormal {
			tags = append(tags, s.buildColumnNameString(field))
			if field.IsPrimaryKey {
				primaryKeys = append(primaryKeys, s.quote(field.DbName))
			}
		}
	}
	var primaryKeyStr string
	if len(primaryKeys) > 0 {
		primaryKeyStr = fmt.Sprintf(", PRIMARY KEY (%v)", strings.Join(primaryKeys, ","))
	}
	var tableOptions string
	if err := s.repo.Exec(fmt.Sprintf("CREATE TABLE %v (%v %v)%s", s.quote(item.TableName), strings.Join(tags, ","), primaryKeyStr, tableOptions)).Error; err != nil {
		glog.Error(err)
	}
}
func (s *MOFSv) quote(str string) string {
	return "`" + str + "`"
}
func (s *MOFSv) updateTable(item md.MDEntity, old md.MDEntity) {
	//更新栏目
	for i := range item.Fields {
		field := item.Fields[i]
		if field.DbName == "-" || field.DbName == "" {
			continue
		}
		oldField := md.MDField{}
		for oi, ov := range old.Fields {
			if strings.ToLower(field.Code) == strings.ToLower(ov.Code) {
				oldField = old.Fields[oi]
				break
			}
		}
		if !field.IsNormal {
			continue
		}
		newString := s.buildColumnNameString(field)
		if s.repo.Dialect().HasColumn(item.TableName, field.DbName) { //字段已存在
			oldString := s.buildColumnNameString(oldField)
			//修改字段类型、类型长度、默认值、注释
			if oldString != newString {
				if err := s.repo.Exec(fmt.Sprintf("ALTER TABLE %v MODIFY %v;", s.quote(item.TableName), newString)).Error; err != nil {
					glog.Error(err)
				}
			}
		} else { //新增字段
			if err := s.repo.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v;", s.quote(item.TableName), newString)).Error; err != nil {
				glog.Error(err)
			}
		}
	}
}
func (s *MOFSv) buildColumnNameString(item md.MDField) string {
	/*
		column_definition:
		data_type [NOT NULL | NULL] [DEFAULT {literal | (expr)} ]
		[AUTO_INCREMENT] [UNIQUE [KEY]] [[PRIMARY] KEY]
		[COMMENT 'string']

	*/
	fieldStr := s.quote(item.DbName)
	if item.IsPrimaryKey && item.TypeID == md.FIELD_TYPE_STRING {
		fieldStr += " NVARCHAR(36)  NOT NULL"
	} else if item.IsPrimaryKey && item.TypeID == md.FIELD_TYPE_INT {
		fieldStr += " BIGINT NOT NULL"
	} else if item.TypeID == "string" {
		if item.Length <= 0 {
			item.Length = 50
		}
		if item.Length >= 8000 {
			fieldStr += " LONGTEXT"
		} else if item.Length >= 4000 {
			fieldStr += " TEXT"
		} else {
			fieldStr += fmt.Sprintf(" NVARCHAR(%d)", item.Length)
		}

		if !item.Nullable {
			fieldStr += " NOT NULL"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == md.FIELD_TYPE_BOOL {
		fieldStr += " TINYINT NOT NULL"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeID == md.FIELD_TYPE_DATE || item.TypeID == md.FIELD_TYPE_DATETIME {
		fieldStr += " TIMESTAMP"
		if !item.Nullable {
			fieldStr += " NOT NULL"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == md.FIELD_TYPE_DECIMAL {
		fieldStr += " DECIMAL(24,9) NOT NULL"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeID == md.FIELD_TYPE_INT {
		if item.Length == 0 {
			fieldStr += " INT NOT NULL"
		} else if item.Length < 2 {
			fieldStr += " TINYINT NOT NULL"
		} else if item.Length > 4 {
			fieldStr += " SMALLINT NOT NULL"
		} else if item.Length < 8 {
			fieldStr += " INT NOT NULL"
		} else if item.Length >= 8 {
			fieldStr += " BIGINT NOT NULL"
		} else {
			fieldStr += " INT NOT NULL"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeType == md.TYPE_ENTITY || item.TypeType == md.TYPE_ENUM {
		fieldStr += " nvarchar(36)"
		if !item.Nullable {
			fieldStr += " NOT NULL"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else {
		if item.Length <= 0 {
			item.Length = 255
		}
		fieldStr += fmt.Sprintf(" nvarchar(%d)", item.Length)
		if !item.Nullable {
			fieldStr += " NOT NULL"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	}
	fieldStr += fmt.Sprintf(" COMMENT '%s'", item.Name)
	return fieldStr

}
func (s *MOFSv) AddMDEntities(items []md.MDEntity) error {
	entityIds := make([]string, 0)
	oldEntities := make([]md.MDEntity, 0)
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := md.MDEntity{}
		if s.repo.Model(oldEntity).Preload("Fields").Order("id").Where("id=?", entity.ID).Take(&oldEntity); oldEntity.ID != "" {
			oldEntities = append(oldEntities, oldEntity)
			datas := make(map[string]interface{})
			if oldEntity.TableName != entity.TableName {
				datas["TableName"] = entity.TableName
			}
			if oldEntity.Type != entity.Type {
				datas["Type"] = entity.Type
			}
			if oldEntity.Code != entity.Code {
				datas["Code"] = entity.Code
			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if oldEntity.System != entity.System {
				datas["System"] = entity.System
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Memo != entity.Memo {
				datas["Memo"] = entity.Memo
			}
			if len(datas) > 0 {
				s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			if entity.Type == md.TYPE_ENTITY && entity.TableName == "" {
				entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
			}
			s.repo.Create(&entity)
		}
		entityIds = append(entityIds, entity.ID)
	}
	//属性字段
	for i, _ := range items {
		entity := items[i]
		if entity.ID != "" && entity.Type == md.TYPE_ENTITY && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, _ := range entity.Fields {
				field := entity.Fields[f]
				itemCodes = append(itemCodes, field.Code)
				field.IsNormal = true
				if field.DbName == "-" {
					field.IsNormal = false
				}
				if field.DbName == "" && field.IsNormal {
					field.DbName = utils.SnakeString(field.Code)
				}
				oldField := md.MDField{}
				if fieldType := s.getEntity(field.TypeID); fieldType.ID != "" {
					field.TypeType = fieldType.Type
					field.TypeID = fieldType.ID
				}
				if field.TypeType == md.TYPE_ENTITY { //实体
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = false
				} else if field.TypeType == md.TYPE_ENUM { //枚举
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = false
				}
				if s.repo.Model(oldField).Order("id").Where("entity_id=? and code=?", entity.ID, field.Code).Take(&oldField); oldField.ID != "" {
					datas := make(map[string]interface{})
					if oldField.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldField.Tags != field.Tags {
						datas["Tags"] = field.Tags
					}
					if oldField.DbName != field.DbName {
						datas["DbName"] = field.DbName
					}
					if oldField.IsNormal != field.IsNormal {
						datas["IsNormal"] = field.IsNormal
					}
					if oldField.IsPrimaryKey != field.IsPrimaryKey {
						datas["IsPrimaryKey"] = field.IsPrimaryKey
					}
					if oldField.Length != field.Length {
						datas["Length"] = field.Length
					}
					if oldField.Nullable != field.Nullable {
						datas["Nullable"] = field.Nullable
					}
					if oldField.DefaultValue != field.DefaultValue {
						datas["DefaultValue"] = field.DefaultValue
					}
					if oldField.TypeID != field.TypeID {
						datas["TypeID"] = field.TypeID
					}
					if oldField.TypeType != field.TypeType {
						datas["TypeType"] = field.TypeType
					}
					if oldField.Limit != field.Limit {
						datas["Limit"] = field.Limit
					}
					if oldField.MinValue != field.MinValue {
						datas["MinValue"] = field.MinValue
					}
					if oldField.MaxValue != field.MaxValue {
						datas["MaxValue"] = field.MaxValue
					}
					if oldField.Precision != field.Precision {
						datas["Precision"] = field.Precision
					}
					if oldField.AssociationKey != field.AssociationKey {
						datas["AssociationKey"] = field.AssociationKey
					}
					if oldField.ForeignKey != field.ForeignKey {
						datas["ForeignKey"] = field.ForeignKey
					}
					if oldField.Kind != field.Kind {
						datas["Kind"] = field.Kind
					}
					if oldField.Sequence != field.Sequence {
						datas["Sequence"] = field.Sequence
					}
					if len(datas) > 0 {
						s.repo.Model(oldField).Where("entity_id=? and code=?", entity.ID, field.Code).Updates(datas)
					}
				} else {
					s.repo.Create(&field)
				}
			}
			s.repo.Delete(md.MDEntity{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
		}
	}
	//枚举
	for _, entity := range items {
		if entity.ID != "" && entity.Type == md.TYPE_ENUM && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, field := range entity.Fields {
				newEnum := md.MDEnum{ID: field.Code, EntityID: entity.ID, Sequence: f, Name: field.Name}
				oldEnum := md.MDEnum{}
				itemCodes = append(itemCodes, newEnum.ID)
				if s.repo.Model(oldEnum).Order("id").Where("entity_id=? and id=?", newEnum.EntityID, newEnum.ID).Take(&oldEnum); oldEnum.ID != "" {
					datas := make(map[string]interface{})
					if oldEnum.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldEnum.Sequence != newEnum.Sequence {
						datas["Sequence"] = newEnum.Sequence
					}
					if len(datas) > 0 {
						s.repo.Model(oldEnum).Where("entity_id=? and id=?", oldEnum.EntityID, oldEnum.ID).Updates(datas)
					}
				} else {
					s.repo.Create(&newEnum)
				}
			}
			s.repo.Delete(md.MDEntity{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
		}
	}
	if len(entityIds) > 0 {
		toTables := make([]md.MDEntity, 0)
		s.repo.Model(md.MDEntity{}).Preload("Fields").Where("id in (?) and type=?", entityIds, md.TYPE_ENTITY).Find(&toTables)
		return s.entityToTables(toTables, oldEntities)
	}
	//缓存
	md.CacheMD(s.repo)
	return nil
}
func (s *MOFSv) AddActionCommand(items []md.MDActionCommand) error {
	for i, _ := range items {
		entity := items[i]
		if entity.Type == "" && entity.Code == "" {
			continue
		}
		oldEntity := md.MDActionCommand{}
		if s.repo.Model(oldEntity).Where("page_id=? and code=? and type=?", entity.PageID, entity.Code, entity.Type).Order("id").Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Path != entity.Path {
				datas["Path"] = entity.Path
			}
			if oldEntity.Content != entity.Content {
				datas["Content"] = entity.Content
			}
			if oldEntity.Rules != entity.Rules {
				datas["Rules"] = entity.Rules
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.System.NotEqual(entity.System) && entity.System.Valid() {
				datas["System"] = entity.System
			}
			if len(datas) > 0 {
				s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			s.repo.Create(&entity)
		}
	}
	return nil
}

func (s *MOFSv) AddActionRule(items []md.MDActionRule) error {
	for i, _ := range items {
		entity := items[i]
		if entity.Code == "" {
			continue
		}
		oldEntity := md.MDActionRule{}
		if s.repo.Model(oldEntity).Where("code=? and domain=?", entity.Code, entity.Domain).Order("id").Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.System.NotEqual(entity.System) && entity.System.Valid() {
				datas["System"] = entity.System
			}
			if oldEntity.Async.NotEqual(entity.System) && entity.Async.Valid() {
				datas["Async"] = entity.Async
			}
			if len(datas) > 0 {
				s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			s.repo.Create(&entity)
		}
	}
	return nil
}
func (s *MOFSv) AddPage(items []md.MDPage) error {
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := md.MDPage{}
		if s.repo.Model(oldEntity).Where("id=?", entity.ID).Order("id").Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Type != entity.Type {
				datas["Type"] = entity.Type
			}
			if oldEntity.Widgets.NotEqual(entity.Widgets) && entity.Widgets.Valid() {
				datas["Widgets"] = entity.Widgets
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
			}
			if oldEntity.MainEntity != entity.MainEntity && entity.MainEntity != "" {
				datas["MainEntity"] = entity.MainEntity
			}
			if oldEntity.System.NotEqual(entity.System) && entity.System.Valid() {
				datas["System"] = entity.System
			}
			if oldEntity.Templated.NotEqual(entity.Templated) && entity.Templated.Valid() {
				datas["Templated"] = entity.Templated
			}
			if len(datas) > 0 {
				if err := s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas).Error; err != nil {
					glog.Error(err)
				}
			}
		} else {
			if err := s.repo.Create(&entity).Error; err != nil {
				glog.Error(err)
			}
		}
	}
	return nil
}

func (s *MOFSv) AddPageView(items []md.MDPageView) error {
	for i, _ := range items {
		entity := items[i]
		if entity.PageID == "" || entity.Code == "" {
			continue
		}
		oldEntity := md.MDPageView{}
		if s.repo.Model(oldEntity).Order("id").Where("page_id=? and code=?", entity.PageID, entity.Code).Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Data.NotEqual(entity.Data) && entity.Data.Valid() {
				datas["Data"] = entity.Data
			}
			if oldEntity.EntityID != entity.EntityID {
				datas["EntityID"] = entity.EntityID
			}
			if oldEntity.Nullable.NotEqual(entity.Nullable) && entity.Nullable.Valid() {
				datas["Nullable"] = entity.Nullable
			}
			if oldEntity.IsMain.NotEqual(entity.IsMain) && entity.IsMain.Valid() {
				datas["IsMain"] = entity.IsMain
			}
			if oldEntity.Multiple.NotEqual(entity.Multiple) && entity.Multiple.Valid() {
				datas["Multiple"] = entity.Multiple
			}
			if oldEntity.PrimaryKey != entity.PrimaryKey {
				datas["PrimaryKey"] = entity.PrimaryKey
			}
			if len(datas) > 0 {
				s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			s.repo.Create(&entity)
		}
	}
	return nil
}
