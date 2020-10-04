package main

import (
	"regexp"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	testOracle()
}
func testOracle() {

	//exp := "([\\S]+)(?i:(?:as|[\\s])+)([\\S]+)"
	//exp := "([\\S]+.*\\S)(?i:\\s+as+\\s)([\\S]+)|([\\S]+.*[\\S]+)"
	exp := "(?i)([\\S]+.*\\S)(?:\\s)(desc|asc)|([\\S]+.*[\\S]+)"
	strList := []string{
		"AAA as a",
		" fieldAs desc ",
		"fieldA+field asc",
		" SUM( fieldA +    +fieldB +sum(fieldC)   )  AS AS ",
		" SUM( fieldA ) ",
	}
	r := regexp.MustCompile(exp)
	for _, str := range strList {
		matched := r.FindStringSubmatch(str)
		glog.Error(matched)
	}

}

type testTable struct {
	ID       string      `gorm:"primary_key;size:50" json:"id"`
	IsSystem int         `gorm:"default:11"`
	Code     string      `gorm:"default:code11"`
	Value    utils.SJson `gorm:"default:1" json:"value"`
	Enabled  utils.SBool `gorm:"default:true;name:启用" json:"enabled"`
}

func (s testTable) MD() *md.Mder {
	return &md.Mder{ID: "test.table", Domain: "test", Name: "标签"}
}

func (s testTable) TableName() string {
	return "testTable"
}
