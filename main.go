package main

import (
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		glog.Errorf("%v\r\n", utils.GUID())
	}
	//test_query()
}
func test_query() {
	mysql := repositories.NewMysqlRepo()

	if err := di.Global.Provide(func() *repositories.MysqlRepo {
		return mysql
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}

	exector:=query.NewExector("amb_elements")
	exector.Where("Purpose.Name like 'a'")
	q,_:=exector.PrepareQuery(mysql)
	glog.Error(q.QueryExpr())

	glog.Errorf("%v\r\n","dd")
}
