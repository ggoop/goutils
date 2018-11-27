package context

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ggoop/goutils/configs"

	"github.com/dgrijalva/jwt-go"
)

type Context struct {
	data map[string]string
}
type Config struct {
	Expiration bool
	Credential bool
	User       bool
	Ent        bool
}

const (
	AuthSessionKey    = "GGOOPAUTH"
	DefaultContextKey = "context"
)

func (s *Context) Copy() *Context {
	c := Context{}
	c.data = make(map[string]string)
	if s.data != nil {
		for k, v := range s.data {
			c.data[k] = v
		}
	}
	return &c
}

func (s *Context) SetEntID(ent string) {
	s.SetValue("ent_id", ent)
}
func (s *Context) SetClientID(ent string) {
	s.SetValue("client_id", ent)
}
func (s *Context) SetUserID(user string) {
	s.SetValue("user_id", user)
}
func (s *Context) EntID() string {
	return s.GetValue("ent_id")
}
func (s *Context) ClientID() string {
	return s.GetValue("client_id")
}
func (s *Context) UserID() string {
	return s.GetValue("user_id")
}
func (s *Context) ExpiresAt() int64 {
	str := s.GetValue("expires_at")
	exp, _ := strconv.ParseInt(str, 10, 64)
	return exp
}

// 验证
func (s *Context) Valid(user, ent bool) error {
	if user && s.VerifyUser() != nil {
		return fmt.Errorf("用户验证失败")
	}
	if ent && s.VerifyEnt() != nil {
		return fmt.Errorf("企业验证失败")
	}
	return nil
}
func (c *Context) VerifyClient() error {
	if c.ClientID() == "" {
		return fmt.Errorf("客户端验证失败")
	}
	return nil
}
func (c *Context) VerifyUser() error {
	if c.UserID() == "" {
		return fmt.Errorf("用户验证失败")
	}
	return nil
}
func (c *Context) VerifyEnt() error {
	if c.EntID() == "" {
		return fmt.Errorf("企业验证失败")
	}
	return nil
}
func (c *Context) VerifyExpiresAt(cmp int64, required bool) bool {
	exp := c.ExpiresAt()
	if exp == 0 {
		return !required
	}
	return cmp <= exp
}

//
func (s *Context) Clean() {
	s.data = make(map[string]string)
}
func (s *Context) GetValue(name string) string {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	return s.data[strings.ToLower(name)]
}
func (s *Context) GetIntValue(name string) int {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	v, _ := strconv.Atoi(s.data[strings.ToLower(name)])
	return v
}
func (s *Context) SetValue(name, value string) *Context {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	name = strings.ToLower(name)
	if name == "ent_id" {
		s.data[name] = value
	} else if name == "user_id" {
		s.data[name] = value
	} else {
		s.data[name] = value
	}
	return s
}

func (s *Context) ToTokenString() string {
	claim := jwt.MapClaims{}
	for k, v := range s.data {
		claim[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(configs.Default.App.Token))
	if err != nil {
		return ""
	}
	return "bearer " + tokenString
}
func (s *Context) FromTokenString(token string) (*Context, error) {
	ctx := Context{}
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		token = tokenParts[1]
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Default.App.Token), nil
	})
	if err != nil {
		return &ctx, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		for k, v := range claims {
			if vstr, ok := v.(string); ok {
				ctx.SetValue(k, vstr)
			}
		}
	}
	return &ctx, nil
}
