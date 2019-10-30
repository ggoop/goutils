package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	//glog.SetPath(utils.JoinCurrentPath(configs.Default.Log.Path))
	//glog.AddLogFile(utils.JoinCurrentPath(configs.Default.Log.Path))

	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	//glog.Error(utils.GetIpAddrs())

	test_query()
}
func test_query() {
	glog.Error(utils.GetIpAddrs())
}
