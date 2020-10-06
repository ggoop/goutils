package main

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

type testTag struct {
	ID     string `gorm:"primary_key;size:50" json:"id"`
	Code   string `gorm:"size:2"`
	Name   utils.SBool
	Json   utils.SJson
	TypeID string     `gorm:"size:10"`
	Type   *md.MDEnum `gorm:"limit:md.type.enum"`
	Field  md.MDField
}

func (s *testTag) MD() *md.Mder {
	return &md.Mder{ID: "test.tag", Domain: "test", Name: "标签"}
}

func TestOQL(t *testing.T) {
	oql := md.GetOQL()
	oql.From("cbo_depts").Select("id as id,code as Code,name")
	oql.Where("id>? and id=?", 1, 2)
	var count int
	parseValues := oql.Count(&count)
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
func TestEnumPreload(t *testing.T) {
	repo := repositories.Default()

	md.Migrate(repo)
	count := 0
	glog.Error(repo.Table("DEPT").Select("DNAME").Where("DNAME=?", "SALES").Count(&count).Error)

	glog.Info(utils.GUID())

}
func TestSplitMatched(t *testing.T) {
	REGEXP_VAR_EXP := `[,|;|，|；\|]`
	str := "a1,b2;c33，e4；f55"
	str = "0|44"
	r, _ := regexp.Compile(REGEXP_VAR_EXP)
	matched_strict := r.Split(str, -1)
	ss := strings.Join(matched_strict, ";")
	t.Error(ss)
}
func TestQuerySJson_Parse(t *testing.T) {
	str := ""
	var jsonData interface{}
	if err := json.Unmarshal([]byte(str), &jsonData); err != nil {
		glog.Error(err)
	}
	glog.Error(jsonData)
}
func TestQueryMigrate(t *testing.T) {
	db := repositories.Default()
	md.InitMD_Completed = true
	md.Migrate(db, &testTag{})

	items := make([]testTag, 0)
	item := testTag{Json: utils.SJson_Parse([]string{"fdsaf", "fdsafddddd"})}
	//
	db.Last(&item)
	db.Find(&items)
	t.Log(item)
	t.Log(items)
}
func TestValueParamMatched(t *testing.T) {
	str := "$$df.i_di> {aaa} AAA} @{a0}@ {a_a} = @ent +ddd +@ent+@entd"
	r, _ := regexp.Compile(query.REGEXP_VAR_EXP)
	matched_strict := r.FindAllStringSubmatch(str, -1)
	t.Error(matched_strict)
}
func TestExectorMatched(t *testing.T) {
	str := "$$df.i_di> {aaa} {AAA} {a0} {a_a}"
	r, _ := regexp.Compile(query.REGEXP_FIELD_EXP_STRICT)
	matched_strict := r.FindAllStringSubmatch(str, -1)
	t.Error(matched_strict)

	r, _ = regexp.Compile(query.REGEXP_FIELD_EXP)
	matched := r.FindAllStringSubmatch(str, -1)
	t.Error(matched)
}
func TestQueryQuery(t *testing.T) {
	db := repositories.Default()
	exector := query.NewExector("test.tag")
	exector.Select("ID").Select("Code")
	//exector.Where("id =? or $$Entid =?", "01ea463a1d1400439ab600505681e808", "01ea24bfd81ec140333bacde48001122")
	qqqqq, err := exector.Query(db)
	if err != nil {
		t.Error(err)
	}
	t.Error(qqqqq)
}
