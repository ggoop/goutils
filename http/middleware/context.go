package middleware

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/context"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type Context struct {
	Sessions *sessions.Sessions
}
type ContextHandle struct {
	Sessions *sessions.Sessions

	Config context.Config
}

func (m *Context) New(config context.Config) func(ctx iris.Context) {
	item := ContextHandle{
		Sessions: m.Sessions,
		Config:   config,
	}
	return item.Handle
}
func (m *Context) Default(ctx iris.Context) {
	item := ContextHandle{
		Sessions: m.Sessions,
	}
	item.Handle(ctx)
}
func (m *ContextHandle) Handle(ctx iris.Context) {
	var (
		uc  *context.Context
		err error
	)
	if IsJWTContext(ctx) {
		if uc, err = m.CheckJWT(ctx); err != nil {
			ctx.StopExecution()
			return
		}
	} else {
		if uc, err = m.CheckSession(ctx); err != nil {
			ctx.StopExecution()
			return
		}
	}
	if uc == nil {
		uc = &context.Context{}
	}
	if uc.EntID() == "" {
		entID := ctx.GetHeader("ENT")
		if entID != "" {
			uc.SetEntID(entID)
		}
	}
	ctx.Values().Set(context.DefaultContextKey, uc)
	if ctx.IsAjax() {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}
	ctx.Next()

}
func (m *ContextHandle) CheckSession(ctx iris.Context) (*context.Context, error) {
	var (
		uc   *context.Context
		err  error
		find bool
	)

	if m.Sessions != nil {
		session := m.Sessions.Start(ctx)
		if v := session.Get(context.AuthSessionKey); v != nil {
			if str, ok := v.(string); ok {
				uc, err = (&context.Context{}).FromTokenString(str)
			} else if obj, ok := v.(*context.Context); ok {
				uc = obj
			}
		}
		find = true
	}
	if !find {
		if str := ctx.GetCookie(context.AuthSessionKey); str != "" {
			uc, err = uc.FromTokenString(str)
			if err == nil {
				find = true
			}
		}
	}
	if !find {
		if m.Config.Credential {
			return uc, fmt.Errorf("Required authorization token not found")
		}
		return uc, nil
	}
	return uc, nil
}
func (m *ContextHandle) CheckJWT(ctx iris.Context) (*context.Context, error) {
	var (
		uc  *context.Context
		err error
	)
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return uc, fmt.Errorf("Authorization header format must be Bearer {token}")
	}
	if uc, err = uc.FromTokenString(token); err != nil {
		return uc, err
	}
	return uc, nil
}

func IsJWTContext(ctx iris.Context) bool {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return false
	}
	return true
}
