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
	di.SetGlobal(di.New())
	di.Global.Provide(repositories.NewMysqlRepo)

	if err := di.Global.Invoke(func(mysql *repositories.MysqlRepo) {
		q := query.QueryCase{}
		q.Query = &query.Query{Entry: "CboTeam", Code: "CboTeam"}
		q.Columns = []query.QueryColumn{
			query.QueryColumn{Field: "Code"},
			query.QueryColumn{Field: "Dept.Org.Name"},
			query.QueryColumn{Field: "Org.Name"},
			query.QueryColumn{Field: "Name"},
		}
		query.RunCase(mysql, q)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
