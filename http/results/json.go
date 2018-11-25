package results

import (
	"fmt"
	"strings"

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
func NotFound(params ...string) Result {
	return ToError(fmt.Errorf("找不到 %s", strings.Join(params, " ")), iris.StatusNotFound)
}
func ToError(err interface{}, code ...int) Result {
	if code != nil && len(code) > 0 {
		if ev, ok := err.(error); ok {
			return mvc.Response{
				Code:   code[0],
				Object: iris.Map{"msg": ev.Error()},
			}
		} else {
			return mvc.Response{
				Code:   code[0],
				Object: iris.Map{"msg": err},
			}
		}

	} else {
		if ev, ok := err.(error); ok {
			return mvc.Response{
				Code:   iris.StatusBadRequest,
				Object: iris.Map{"msg": ev.Error()},
			}
		} else {
			return mvc.Response{
				Code:   iris.StatusBadRequest,
				Object: iris.Map{"msg": err},
			}
		}
	}
}
