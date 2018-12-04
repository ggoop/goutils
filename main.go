package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func main() {
	claims := make(map[string]interface{})
	claims["aaa"] = 1111
	claims["exp"] = time.Now().Add((-time.Hour) * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()

	str := utils.CreateJWTToken(claims, "aaa")

	glog.Errorf("sss ,%v ", str)

	mapClaims, err := utils.ParseJWTToken(str, "aaa")
	licValues := make([]string, 0)
	for k, v := range mapClaims {
		licValues = append(licValues, fmt.Sprintf("%v:%v", k, v))
	}
	if mapClaims.VerifyExpiresAt(time.Now().Unix(), false) {
		glog.Errorf("sss ,%v ", str)
	}
	glog.Errorf("sss ,%v %v %v", mapClaims, err, strings.Join(licValues, ","))

}
