package main

import (
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

}

type testTag struct {
	ID   string `gorm:"primary_key;size:50" json:"id"`
	Code string `gorm:"size:2"`
	Name md.SBool
}

func (s *testTag) MD() *md.Mder {
	return &md.Mder{ID: "test.tag", Domain: "test", Name: "标签"}
}
