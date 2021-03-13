package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/http/middleware"
)

// 路由服务注册
func Register() {
	registerView()
}
func registerView() {
	if err := di.Global.Invoke(func(app *iris.Application, contextMid *middleware.Context) {
		{
			m := mvc.New(app.Party("/api/md", contextMid.Default))
			m.Handle(new(MdController))
		}
		app.Options("{directory:path}", func(ctx iris.Context) {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Headers", "*")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Write([]byte("ok"))
		})
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
