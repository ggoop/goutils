package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	startHttp()
}
