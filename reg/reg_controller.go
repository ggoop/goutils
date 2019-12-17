package reg

import (
	"github.com/ggoop/goutils/http/results"
	"github.com/kataras/iris"
)

type RegController struct {
	Ctx   iris.Context
	Store *RegStoreSv
}

func (c *RegController) PostRegister() results.Result {
	item := RegObject{}
	if err := c.Ctx.ReadJSON(&item); err != nil {
		return results.ToError(err)
	}
	c.Store.Add(item)
	return results.ToJson(results.Map{"data": true})
}
func (c *RegController) GetBy(code string) results.Result {
	return results.ToJson(results.Map{"data": c.Store.Get(RegObject{Code: code})})
}
func (c *RegController) Get() results.Result {
	return results.ToJson(results.Map{"data": c.Store.GetAll()})
}
