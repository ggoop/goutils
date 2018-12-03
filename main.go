package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	str, err := utils.Encrypt("123", "aaa")
	if err != nil {
		glog.Errorf("sss ,%v %v", str, err)
	}

	str, err = utils.Decrypt(str, "aaa")
	if err != nil {
		glog.Errorf("sss ,%v %v", str, err)
	}
}
