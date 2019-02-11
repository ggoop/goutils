package results

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type (
	Result = mvc.Result
	Map    = iris.Map
)

func ToJson(data interface{}) Result {
	return mvc.Response{
		Object: data,
	}
}
func Unauthenticated() Result {
	return ToError(fmt.Errorf("Unauthenticated"), iris.StatusUnauthorized)
}
func ParamsRequired(params ...string) Result {
	return ToError(fmt.Errorf("参数 %s 不能为空!", params), iris.StatusBadRequest)
}
func ParamsFailed(params ...string) Result {
	return ToError(fmt.Errorf("参数 %s 不正确!", params), iris.StatusUnsupportedMediaType)
}
func NotFound(params ...string) Result {
	return ToError(fmt.Errorf("找不到 %s", strings.Join(params, " ")), iris.StatusNotFound)
}
func ToSingle(data interface{}) Result {
	return mvc.Response{
		Object: iris.Map{"data": data},
	}
}
func ToError(err interface{}, code ...int) Result {
	res := mvc.Response{}
	obj := iris.Map{}
	if ev, ok := err.(utils.GError); ok {
		obj["msg"] = ev.Error()
	} else if ev, ok := err.(error); ok {
		obj["msg"] = ev.Error()
	} else {
		obj["msg"] = err
	}
	if code != nil && len(code) > 0 {
		if code[0] >= 100 && code[0] < 1000 {
			res.Code = code[0]
		} else {
			obj["code"] = code[0]
			res.Code = iris.StatusBadRequest
		}
	} else {
		res.Code = iris.StatusBadRequest
	}
	if code != nil && len(code) > 1 {
		obj["code"] = code[1]
	}
	res.Object = obj
	return res
}
