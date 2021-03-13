package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/kataras/iris/sessions"
	"time"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/http/middleware"
	"github.com/ggoop/goutils/http/middleware/logger"
	"github.com/ggoop/goutils/http/middleware/recover"
	routes "github.com/ggoop/goutils/httpApp"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func startHttp() error {
	app := iris.New()
	// 创建容器
	di.SetGlobal(di.New())

	app.Use(recover.New())
	app.Use(logger.New())

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{"msg": "Not Found : " + ctx.Path()})
	})

	// 启用session
	sessionManager := sessions.New(sessions.Config{
		Cookie:       "site_session_id",
		Expires:      48 * time.Hour,
		AllowReclaim: true,
	})
	hero.Register(sessionManager.Start)
	if err := di.Global.Provide(func() *sessions.Sessions {
		return sessionManager
	}); err != nil {
		return glog.Errorf("注册缓存服务异常:%s", err)
	}
	// 注册app
	if err := di.Global.Provide(func() *iris.Application {
		return app
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	//注册中间件
	if err := di.Global.Provide(func() *middleware.Context {
		return &middleware.Context{Sessions: sessionManager}
	}); err != nil {
		return glog.Errorf("注册上下文中间件异常:%s", err)
	}
	repo := repositories.NewMysqlRepo()
	md.Migrate(repo)
	// 路由注册
	routes.Register()
	// 启动服务
	dc := iris.DefaultConfiguration()
	if utils.PathExists("env/iris.yaml") {
		dc = iris.YAML(utils.JoinCurrentPath("env/iris.yaml"))
	}
	dc.DisableBodyConsumptionOnUnmarshal = true
	dc.FireMethodNotAllowed = true
	if err := app.Run(iris.Addr(":"+configs.Default.App.Port), iris.WithConfiguration(dc)); err != nil {
		return glog.Errorf("Run service error:%s\n", err.Error())
	}
	return nil
}
