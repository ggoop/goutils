package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	str, err := utils.AesCFBEncrypt("1", "aaa")
	if err != nil {
		glog.Errorf("sss ,%v %v", str, err)
	}

	str, err = utils.AesCFBDecrypt(str, "aaa2")
	if err != nil {
		glog.Errorf("sss ,%v %v", str, err)
	}
}
