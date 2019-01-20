package main

import (
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		glog.Errorf("%v\r\n", utils.GUID())
	}
	// test_query()
}
func test_query() {
	di.SetGlobal(di.New())
	di.Global.Provide(repositories.NewMysqlRepo)

	if err := di.Global.Invoke(func(mysql *repositories.MysqlRepo) {
		md.Migrate(mysql, query.QueryCase{})
		q := query.QueryCase{}
		q.Query = &query.Query{Entry: "MDEntity as a", Code: "MDEntity"}
		q.Columns = []query.QueryColumn{
			query.QueryColumn{Field: "Name"},
			query.QueryColumn{Field: "Type.Name"},
		}
		exector := q.GetExector()
		if qc, err := exector.Query(mysql); err != nil {
			glog.Errorf("errir is $v", err)
		} else {
			m := query.Query{}
			if err := qc.Find(&m).Error; err != nil {
				glog.Errorf("sss is %v", err)
			}
			glog.Errorf("sss is %v", m)
		}
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
