package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}

	//test_query()
}
func test_query() {
	md.Migrate(repositories.NewMysqlRepo(), &md.MDEntity{})
	glog.Error(utils.GetIpAddrs())

}
