package main

import (
	"regexp"
	"strings"

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
	exp := `(?i)([left|join|right|union]+)([A-Za-z0-9._])(?:[as|\s])?([A-Za-z0-9_])?`
	strList := []string{
		"tableA a on aaa=dfds and dfd=ddd",
		"tableB A on fdfd=dfd and ddd=sss",
	}
	r := regexp.MustCompile(exp)
	for _, str := range strList {
		matched := r.FindStringSubmatch(str)
		matchStr := strings.Join(matched, "||")
		glog.Error(matchStr)
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
