package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		glog.Errorf("%v \r\n", utils.GUID())
	}
}
