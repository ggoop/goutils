package main

import (
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
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
		
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
