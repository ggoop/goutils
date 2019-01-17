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
	test_query()
}
func test_query() {
	di.SetGlobal(di.New())
	di.Global.Provide(repositories.NewMysqlRepo)

	if err := di.Global.Invoke(func(mysql *repositories.MysqlRepo) {
		md.Migrate(mysql, query.QueryCase{})
		q := query.QueryCase{}
		q.Query = &query.Query{Entry: "CboTeam as a", Code: "CboTeam"}
		q.Columns = []query.QueryColumn{
			query.QueryColumn{Field: "a.Code"},
			query.QueryColumn{Field: "Dept.Org.Name"},
			query.QueryColumn{Field: "a.Org.Name"},
			query.QueryColumn{Field: "Name"},
		}
		q.Wheres = []query.QueryWhere{
			query.QueryWhere{
				Field: "a.Code", Operator: "=", Value: "1", Children: []query.QueryWhere{
					query.QueryWhere{Field: "name", Operator: "=", Value: "12na'me"},
					query.QueryWhere{Field: "code", Operator: "=", Value: "code"},
				},
			},
			query.QueryWhere{Logical: "or", Children: []query.QueryWhere{
				query.QueryWhere{Field: "id", Operator: "=", Value: "1"},
				query.QueryWhere{Field: "id", Operator: "=", Value: "2"},
			}},
		}
		exector := query.NewCaseExector(q)
		if qc, err := exector.PrepareQuery(mysql); err != nil {
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
