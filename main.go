package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	//test_query()
}
func test_query() {
	mysql:=repositories.NewMysqlRepo()
	md.Migrate(mysql, &TestTag{})

items:=make([]interface{},0)
item:=TestTag{ID:utils.GUID(),TypeID:"dd"}

	item.CreatedAt=md.NewTime()
	items = append(items, item)
	//mysql.Create(item)
	mysql.BatchInsert(items)

}

type TestTag struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt md.Time    `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *md.Time   `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:50"`
	TypeID    string     `gorm:"size:50"`
	Type      *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:cbo.biz.type"`
	ParentID  string     `gorm:"size:50;name:父模型" json:"parent_id"`
}

func (s *TestTag) MD() *md.Mder {
	return &md.Mder{ID: "test.tag", Domain: "test", Name: "标签"}
}
