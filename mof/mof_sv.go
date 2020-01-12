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
func (s *MOFSv) EntityToTables(items ...md.MDEntity) error {
	if items == nil || len(items) == 0 {
		return nil
	}
	for i, _ := range items {
		entity := items[i]
		if entity.Type != md.TYPE_ENTITY {
			continue
		}
		if entity.TableName == "" {
			entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
		}
		if s.repo.Dialect().HasTable(entity.TableName) {
			s.updateTable(entity)
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
func (s *MOFSv) updateTable(item md.MDEntity) {
	//更新栏目
	for i := range item.Fields {
		field := item.Fields[i]
		if field.DbName == "-" {
			continue
		}
		if !field.IsNormal {
			continue
		}
		if s.repo.Dialect().HasColumn(item.TableName, field.DbName) {
			continue
		}
		if err := s.repo.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v;", s.quote(item.TableName), s.buildColumnNameString(field))).Error; err != nil {
			glog.Error(err)
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
	if item.IsPrimaryKey && item.TypeID == "string" {
		fieldStr += " nvarchar(36)  not null"
	} else if item.IsPrimaryKey && item.TypeID == "int" {
		fieldStr += " bigint not null"
	} else if item.TypeID == "string" {
		if item.Length <= 0 {
			item.Length = 50
		}
		if item.Length >= 8000 {
			fieldStr += " LONGTEXT"
		} else if item.Length >= 4000 {
			fieldStr += " TEXT"
		} else {
			fieldStr += fmt.Sprintf(" nvarchar(%d)", item.Length)
		}

		if !item.Nullable {
			fieldStr += " not null"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == "boolean" {
		fieldStr += " bit not null"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeID == "datetime" {
		fieldStr += " TIMESTAMP"
		if !item.Nullable {
			fieldStr += " not null"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == "datetime" {
		fieldStr += " TIMESTAMP"
		if !item.Nullable {
			fieldStr += " not null"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == "decimal" {
		fieldStr += " DECIMAL(24,9) not null"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeID == "int" {
		if item.Length == 0 {
			fieldStr += " int not null"
		} else if item.Length < 2 {
			fieldStr += " tinyint not null"
		} else if item.Length > 4 {
			fieldStr += " smallint not null"
		} else if item.Length < 8 {
			fieldStr += " int not null"
		} else if item.Length >= 8 {
			fieldStr += " bigint not null"
		} else {
			fieldStr += " int not null"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
	} else if item.TypeType == md.TYPE_ENTITY || item.TypeType == md.TYPE_ENUM {
		fieldStr += " nvarchar(36)"
		if !item.Nullable {
			fieldStr += " not null"
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
			fieldStr += " not null"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	}
	return fieldStr

}
func (s *MOFSv) AddMDEntities(items []md.MDEntity) error {
	entityIds := make([]string, 0)
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := md.MDEntity{}
		if s.repo.Model(oldEntity).Order("id").Where("id=?", entity.ID).Take(&oldEntity); oldEntity.ID != "" {
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

		if len(entity.Fields) > 0 {
			for f, _ := range entity.Fields {
				field := entity.Fields[f]
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
				}
				if field.TypeType == md.TYPE_ENTITY {
					if field.ForeignKey == "" {
						field.ForeignKey = field.Code
					}
					if field.AssociationKey == "" {
						field.ForeignKey = "ID"
					}
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
		}
		entityIds = append(entityIds, entity.ID)
	}
	if len(entityIds) > 0 {
		toTables := make([]md.MDEntity, 0)
		s.repo.Model(md.MDEntity{}).Preload("Fields").Where("id in (?) and type=?", entityIds, md.TYPE_ENTITY).Find(&toTables)
		return s.EntityToTables(toTables...)
	}
	return nil
}
func (s *MOFSv) AddActionCommand(items []md.MDActionCommand) error {
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := md.MDActionCommand{}
		if s.repo.Model(oldEntity).Where("id=?", entity.ID).Order("id").Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Path != entity.Path {
				datas["Path"] = entity.Path
			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
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
		if entity.ID == "" {
			continue
		}
		oldEntity := md.MDActionRule{}
		if s.repo.Model(oldEntity).Where("id=?", entity.ID).Order("id").Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Content != entity.Content && entity.Content != "" {
				if entity.Content == "-" {
					datas["Content"] = ""
				} else {
					datas["Content"] = entity.Content
				}

			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
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

func (s *MOFSv) AddActionMaker(items []md.MDActionMaker) error {
	for i, _ := range items {
		entity := items[i]
		if entity.MakerID == "" || entity.MakerType == "" || entity.CommandID == "" || entity.RuleID == "" {
			continue
		}
		oldEntity := md.MDActionMaker{}
		if s.repo.Model(oldEntity).Order("id").Where("maker_id=? and maker_type=? and command_id=? and rule_id=?", entity.MakerID, entity.MakerType, entity.CommandID, entity.RuleID).Take(&oldEntity); oldEntity.ID != "" {
			datas := make(map[string]interface{})
			if oldEntity.Group != entity.Group {
				datas["Group"] = entity.Group
			}
			if oldEntity.Sequence != entity.Sequence {
				datas["Sequence"] = entity.Sequence
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
			if oldEntity.Widgets != entity.Widgets && entity.Widgets != "" {
				if entity.Widgets == "-" {
					datas["Widgets"] = ""
				} else {
					datas["Widgets"] = entity.Widgets
				}
			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
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
			if oldEntity.Data != entity.Data && entity.Data != "" {
				if entity.Data == "-" {
					datas["Data"] = ""
				} else {
					datas["Data"] = entity.Data
				}

			}
			if oldEntity.EntityID != entity.EntityID {
				datas["EntityID"] = entity.EntityID
			}
			if oldEntity.Nullable != entity.Nullable {
				datas["Nullable"] = entity.Nullable
			}
			if oldEntity.IsMain != entity.IsMain {
				datas["IsMain"] = entity.IsMain
			}
			if oldEntity.Multiple != entity.Multiple {
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
