package mof

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"
)

/**
请求通用类
*/
type ReqContext struct {
	PageID     string      `json:"page_id"` //页面ID
	ID         string      `json:"id"`
	IDS        []string    `json:"ids"`
	UserID     string      `json:"user_id"` //用户ID
	EntID      string      `json:"ent_id"`  //企业ID
	OrgID      string      `json:"org_id"`  //组织ID
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Command    string      `json:"command"` // 动作编码
	Rule       string      `json:"rule"`    //规则编码
	URI        string      `json:"uri"`
	Method     string      `json:"method"`
	Q          string      `json:"q"`
	Condition  interface{} `json:"condition"`
	MainEntity string      `json:"main_entity"`
	Data       interface{} `json:"data"` //数据
}
type ResContext struct {
	Data  utils.Map
	Error error
}

func (s *ResContext) SetData(name string, value interface{}) *ResContext {
	if s.Data == nil {
		s.Data = utils.Map{}
	}
	s.Data[name] = value
	return s
}

type PageViewDTO struct {
	Data       interface{} `json:"data"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	EntityID   string      `json:"entity_id"`
	PrimaryKey string      `json:"primary_key"`
	Multiple   md.SBool    `json:"multiple"`
	Nullable   md.SBool    `json:"nullable"`
	IsMain     md.SBool    `json:"is_main"`
}

func (s ReqContext) Copy() ReqContext {
	return ReqContext{
		ID: s.ID, IDS: s.IDS, UserID: s.UserID, EntID: s.EntID, OrgID: s.OrgID,
		PageID: s.PageID,
		Page:   s.Page, PageSize: s.PageSize, Q: s.Q,
		Command: s.Command, Rule: s.Rule,
		URI: s.URI, Method: s.Method,
		Condition: s.Condition, MainEntity: s.MainEntity, Data: s.Data,
	}
}

/**
规则通用接口
*/
type IActionRule interface {
	Exec(context *ReqContext, res *ResContext) error
	GetRule() RuleRegister
}

var _action_rule = initActionRule()

//基础规则
type SimpleRule struct {
	Rule   RuleRegister
	Handle func(context *ReqContext, res *ResContext) error
}

func (s *SimpleRule) Exec(context *ReqContext, res *ResContext) error {
	if s.Handle != nil {
		return s.Handle(context, res)
	}
	return glog.Error("没有实现")
}

func (s *SimpleRule) GetRule() RuleRegister {
	return s.Rule
}

func initActionRule() map[string]IActionRule {
	return make(map[string]IActionRule)
}

func GetActionRule(domain, rule string) (IActionRule, bool) {
	key := fmt.Sprintf("%s:%s", strings.ToLower(domain), strings.ToLower(rule))
	if r, ok := _action_rule[key]; ok {
		return r, ok
	}
	return nil, false
}
func RegisterActionRule(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			key := fmt.Sprintf("%s:%s", strings.ToLower(rule.GetRule().Domain), strings.ToLower(rule.GetRule().Code))
			_action_rule[key] = rule
		}
	}

}

// 注册器
type RuleRegister struct {
	Code   string
	Domain string
}
