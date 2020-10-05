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

	oql := md.GetOQL().From("cbo_depts as a").Select("id as id,code as code,name")
	oql.Where("id>?", 1)
	parseValues := oql.Parse()
	glog.Error(parseValues)
	//exp := "([\\S]+)(?i:(?:as|[\\s])+)([\\S]+)"
	//exp := "([\\S]+.*\\S)(?i:\\s+as+\\s)([\\S]+)|([\\S]+.*[\\S]+)"
	exp := `\(.*,`
	strList := []string{
		"fieldA,fieldB",
		"sum(fieldA,fieldB)",
	}

	r := regexp.MustCompile(exp)
	for _, str := range strList {
		matched := r.MatchString(str)
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
