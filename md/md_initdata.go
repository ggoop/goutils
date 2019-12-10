package md

import (
	"github.com/ggoop/goutils/repositories"
)
func initData(db *repositories.MysqlRepo) {
	items := make([]interface{}, 0)
	itemFields := make([]interface{}, 0)
	enumID := ""
	has := 0
	if db.Model(MDEntity{}).Where("id=?", "md.type.enum").Count(&has); has > 0 {
		return
	}
	//基础数据类型
	items = append(items, &MDEntity{ID: "string", Name: "字符", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "int", Name: "整数", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "boolean", Name: "布尔", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "decimal", Name: "浮点数", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "text", Name: "文本", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "date", Name: "日期", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "datetime", Name: "时间", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "binary", Name: "二进制", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "xml", Name: "XML", Type: TYPE_SIMPLE, Domain: md_domain, System: true})
	items = append(items, &MDEntity{ID: "json", Name: "JSON", Type: TYPE_SIMPLE, Domain: md_domain, System: true})

	//基础枚举
	enumID = "md.type.enum"
	items = append(items, &MDEntity{ID: enumID, Name: "元数据类型", Type: TYPE_ENUM, Domain: md_domain, System: true})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_SIMPLE, Name: "简单类型"})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_ENUM, Name: "枚举"})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_DTO, Name: "对象"})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_ENTITY, Name: "实体"})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_INTERFACE, Name: "接口"})
	itemFields = append(itemFields, &MDField{EntityID: enumID, Code: TYPE_VIEW, Name: "视图"})

	db.BatchInsert(items)
	db.BatchInsert(itemFields)
}
