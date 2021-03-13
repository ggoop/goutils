package routes

import (
	"github.com/kataras/iris"

	"github.com/ggoop/goutils/context"
	"github.com/ggoop/goutils/files"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/http/results"
	"github.com/ggoop/goutils/mof"
)

type MdController struct {
	Ctx iris.Context
}

func (c *MdController) PostUiImport() results.Result {
	ctx := c.Ctx.Values().Get(context.DefaultContextKey).(*context.Context)
	if ctx.IsValid() {

	}
	var postInput mof.ReqContext
	if err := c.Ctx.ReadForm(&postInput); err != nil {
		c.Ctx.ReadJSON(&postInput)
	}
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		return results.ToError(err)
	}

	glog.Error(info.Filename)
	if file != nil {
		if datas, err := files.NewExcelSv().GetExcelDatasByReader(file); err != nil {
			return results.ToError(err)
		} else {
			glog.Error(len(datas))
		}
	}
	return results.ToJson(true)
}
