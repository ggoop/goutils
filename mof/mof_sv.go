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
		if field.IsNormal.IsTrue() {
			tags = append(tags, s.buildColumnNameString(field))
			if field.IsPrimaryKey.IsTrue() {
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
	return s.repo.Dialect().Quote(str)
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
		if !field.IsNormal.IsTrue() {
			continue
		}
		newString := s.buildColumnNameString(field)
		if s.repo.Dialect().HasColumn(item.TableName, field.DbName) { //字段已存在
			oldString := s.buildColumnNameString(oldField)
			//修改字段类型、类型长度、默认值、注释
			if oldString != newString && strings.Contains(item.Tags, "update") {
				if err := s.repo.Exec(fmt.Sprintf("ALTER TABLE %v MODIFY %v", s.quote(item.TableName), newString)).Error; err != nil {
					glog.Error(err)
				}
			}
		} else { //新增字段
			if err := s.repo.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v", s.quote(item.TableName), newString)).Error; err != nil {
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
	dialectName := s.repo.Dialect().GetName()
	if dialectName == "godror" || dialectName == "oracle" {
		return s.buildColumnNameString4Oracle(item)
	} else {
		return s.buildColumnNameString4Mysql(item)
	}
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
				field.IsNormal = utils.SBool_True
				if field.DbName == "-" {
					field.IsNormal = utils.SBool_False
				}
				if field.DbName == "" && field.IsNormal.IsTrue() {
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
					field.IsNormal = utils.SBool_False
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
					field.IsNormal = utils.SBool_False
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
			s.repo.Delete(md.MDField{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
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
			s.repo.Delete(md.MDField{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
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
			if oldEntity.Type != entity.Type && entity.Type != "" {
				datas["Type"] = entity.Type
			}
			if oldEntity.Widgets.NotEqual(entity.Widgets) && entity.Widgets.Valid() {
				datas["Widgets"] = entity.Widgets
			}
			if oldEntity.Extras.NotEqual(entity.Extras) && entity.Extras.Valid() {
				datas["Extras"] = entity.Extras
			}
			if oldEntity.Code != entity.Code && entity.Code != "" {
				datas["Code"] = entity.Code
			}
			if oldEntity.Name != entity.Name && entity.Name != "" {
				datas["Name"] = entity.Name
			}
			if oldEntity.Domain != entity.Domain && entity.Domain != "" {
				datas["Domain"] = entity.Domain
			}
			if oldEntity.MainEntity != entity.MainEntity && entity.MainEntity != "" {
				datas["MainEntity"] = entity.MainEntity
			}
			if oldEntity.Element != entity.Element && entity.Element != "" {
				datas["Element"] = entity.Element
			}
			if oldEntity.System.NotEqual(entity.System) && entity.System.Valid() {
				datas["System"] = entity.System
			}
			if len(datas) > 0 {
				if err := s.repo.Model(oldEntity).Where("id=?", oldEntity.ID).Updates(datas).Error; err != nil {
					return glog.Error(err)
				}
			}
		} else {
			if err := s.repo.Create(&entity).Error; err != nil {
				return glog.Error(err)
			}
		}
		if err := s.savePageWidgets(entity, entity.Children); err != nil {
			return err
		}
	}
	return nil
}

func (s *MOFSv) savePageWidgets(page md.MDPage, widgets []md.MDPageWidget) error {
	for i, _ := range widgets {
		item := widgets[i]
		item.PageID = page.ID
		item.EntID = page.EntID
		item.Sequence = i + 1

		old := md.MDPageWidget{}
		s.repo.Model(old).Where("page_id=? and code=?", item.PageID, item.Code).Order("id").Take(&old)
		if old.ID != "" {
			item.ID = old.ID
			datas := make(map[string]interface{})
			if old.ParentCode != item.ParentCode {
				datas["ParentCode"] = item.ParentCode
			}
			if old.Element != item.Element && item.Element != "" {
				datas["Element"] = item.Element
			}
			if old.Code != item.Code && item.Code != "" {
				datas["Code"] = item.Code
			}
			if old.Name != item.Name && item.Name != "" {
				datas["Name"] = item.Name
			}
			if old.Entity != item.Entity && item.Entity != "" {
				datas["Entity"] = item.Entity
			}
			if old.Field != item.Field && item.Field != "" {
				datas["Field"] = item.Field
			}
			if old.Extras.NotEqual(item.Extras) && item.Extras.Valid() {
				datas["Extras"] = item.Extras
			}
			if old.Required.NotEqual(item.Required) && item.Required.Valid() {
				datas["Required"] = item.Required
			}
			if old.Hidden.NotEqual(item.Hidden) && item.Hidden.Valid() {
				datas["Hidden"] = item.Hidden
			}
			if old.Editable.NotEqual(item.Editable) && item.Editable.Valid() {
				datas["Editable"] = item.Editable
			}
			if old.Placeholder != item.Placeholder && item.Placeholder != "" {
				datas["Placeholder"] = item.Placeholder
			}
			if old.Sequence != item.Sequence && item.Sequence > 0 {
				datas["Sequence"] = item.Sequence
			}
			if old.Value.NotEqual(item.Value) && item.Value.Valid() {
				datas["Value"] = item.Value
			}
			if old.Align != item.Align && item.Align != "" {
				datas["Align"] = item.Align
			}
			if old.Width != item.Width && item.Width != "" {
				datas["Width"] = item.Width
			}
			if old.InputType != item.InputType && item.InputType != "" {
				datas["InputType"] = item.InputType
			}
			if old.DataSource != item.DataSource && item.DataSource != "" {
				datas["DataSource"] = item.DataSource
			}
			if old.DataType != item.DataType && item.DataType != "" {
				datas["DataType"] = item.DataType
			}
			if len(datas) > 0 {
				if err := s.repo.Model(item).Where("id=?", item.ID).Updates(datas).Error; err != nil {
					return glog.Error(err)
				}
			}
		} else {
			item.ID = utils.GUID()
			if err := s.repo.Create(&item).Error; err != nil {
				return glog.Error(err)
			}
		}
	}
	return nil
}
