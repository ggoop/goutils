package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		glog.Errorf("%v\r\n", utils.GUID())
	}
	test_query()
}
func test_query() {
	mysql := repositories.NewMysqlRepo()

	exector:=query.NewExector("sys_oss_objects")
	exector.Where("Oss.Ent.Name like 'a'")
	q,_:=exector.PrepareQuery(mysql)
	glog.Error(q.QueryExpr())

	glog.Errorf("%v\r\n","dd")
}
