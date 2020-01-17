package main

import (
	"net/url"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	//test_query()
}
func test_query() {
	remoteUrl, _ := url.Parse("http://129.9.9.9/aaa")
	remoteUrl.Path = "/api/ents/register"

	glog.Error(remoteUrl.String())
}

type testTag struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt md.Time    `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *md.Time   `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:50"`
	TypeID    string     `gorm:"size:50"`
	Type      *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:cbo.biz.type"`
	ParentID  string     `gorm:"size:50;name:父模型" json:"parent_id"`
}

func (s *testTag) MD() *md.Mder {
	return &md.Mder{ID: "test.tag", Domain: "test", Name: "标签"}
}
