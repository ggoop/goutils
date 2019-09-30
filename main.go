package main

import (
	"github.com/ggoop/goutils/files"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	//glog.SetPath(utils.JoinCurrentPath(configs.Default.Log.Path))
	//glog.AddLogFile(utils.JoinCurrentPath(configs.Default.Log.Path))
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}

	//test_query()
}
func test_query() {
	excor := files.NewExcelSv()
	if data, err := excor.GetExcelData("/Users/samw/project/suite/suite-docs/func.xlsx"); err != nil {
		glog.Error(err)
	} else {
		glog.Error(data)
	}
}
