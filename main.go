package main

import (
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	//glog.SetPath(utils.JoinCurrentPath(configs.Default.Log.Path))
	//glog.AddLogFile(utils.JoinCurrentPath(configs.Default.Log.Path))
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID(), glog.String("aa", "dd"), glog.Int("a234", 33))
		glog.Infof("Failed to fetch URL: %v", i)
	}
	test_query()
}
func test_query() {
	mysql := repositories.NewMysqlRepo()
	mysql.SetLogger(glog.GetLogger("sql"))
	if err := di.Global.Provide(func() *repositories.MysqlRepo {
		return mysql
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	items:=make([]query.Query,0)
	mysql.Model(query.Query{}).Where("id=?","333").Find(&items)
}
