package mof

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

/**
请求通用类
*/
type ReqContext struct {
	PageID     string                 `json:"page_id" form:"page_id"` //页面ID
	Params     map[string]interface{} `json:"params"  form:"params"`  //一般指 页面 URI 参数
	ID         string                 `json:"id" form:"id"`
	IDS        []string               `json:"ids" form:"ids"`
	UserID     string                 `json:"user_id" form:"user_id"` //用户ID
	EntID      string                 `json:"ent_id" form:"ent_id"`   //企业ID
	OrgID      string                 `json:"org_id" form:"org_id"`   //组织ID
	Page       int                    `json:"page" form:"page"`
	PageSize   int                    `json:"page_size" form:"page_size"`
	Command    string                 `json:"command" form:"command"` // 动作编码
	Rule       string                 `json:"rule" form:"rule"`       //规则编码
	URI        string                 `json:"uri" form:"uri"`
	Method     string                 `json:"method" form:"method"`
	Q          string                 `json:"q" form:"q"`                 //模糊查询条件
	Condition  interface{}            `json:"condition" form:"condition"` //附件条件
	MainEntity string                 `json:"main_entity" form:"main_entity"`
	Data       interface{}            `json:"data" form:"data"` //数据
	Tag        string                 `json:"tag" form:"tag"`
	Orders     []query.QueryOrder     `json:"orders" form:"orders"`
	Wheres     []query.QueryWhere     `json:"wheres" form:"wheres"`
	Columns    []query.QueryColumn    `json:"columns" form:"columns"`
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

func (s ResContext) New() ResContext {
	return ResContext{}
}

type PageViewDTO struct {
	Data       interface{} `json:"data"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	EntityID   string      `json:"entity_id"`
	PrimaryKey string      `json:"primary_key"`
	Multiple   utils.SBool `json:"multiple"`
	Nullable   utils.SBool `json:"nullable"`
	IsMain     utils.SBool `json:"is_main"`
}

func (s ReqContext) New() ReqContext {
	return ReqContext{
		UserID: s.UserID, EntID: s.EntID, OrgID: s.OrgID,
	}
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

//执行命令
func DoAction(req ReqContext) (interface{}, error) {
	if req.MainEntity == "" && req.PageID != "" {
		pageMD := md.MDPage{}
		if err := repositories.Default().Model(pageMD).Where("id=?", req.PageID).Take(&pageMD).Error; err != nil {
			return nil, err
		}
		req.MainEntity = pageMD.MainEntity
	}
	res := &ResContext{}
	rules := make([]IActionRule, 0)
	comName := "common"
	//优先指定规则
	if req.Rule != "" {
		ruleCodes := strings.Split(req.Rule, ";")
		for _, r := range ruleCodes {
			rule, ok := GetActionRule(req.PageID, r)
			//没有找到，则查找公共规则
			if !ok && req.PageID != comName {
				rule, ok = GetActionRule(comName, r)
			}
			if rule != nil && ok {
				rules = append(rules, rule)
			}
		}
	} else {
		command := md.MDActionCommand{}
		//查询页面命令
		repositories.Default().Where("page_id=? and code=? and type=?", req.PageID, req.Command, "action").Take(&command)
		if command.ID == "" { //查找实体命令
			repositories.Default().Where("page_id=? and code=? and type=?", req.MainEntity, req.Command, "entity").Take(&command)
		}
		if command.ID == "" { //查找公共命令
			repositories.Default().Where("page_id=? and code=? and type=?", comName, req.Command, "action").Take(&command)
		}
		if command.ID != "" && command.Rules != "" {
			ruleCodes := strings.Split(command.Rules, ";")
			for _, r := range ruleCodes {
				rule, ok := GetActionRule(command.PageID, r)
				//没有找到，则查找公共规则
				if !ok && command.PageID != comName {
					rule, ok = GetActionRule(comName, r)
				}
				if rule != nil && ok {
					rules = append(rules, rule)
				}
			}
		}
	}
	//如果没有找到任何规则，则使用和命令同名规则
	if len(rules) == 0 {
		//页面级
		rule, ok := GetActionRule(req.PageID, req.Command)
		if !ok && req.PageID != req.MainEntity { //实体级
			rule, ok = GetActionRule(req.MainEntity, req.Command)
		}
		if !ok && req.PageID != comName { //公共级
			rule, ok = GetActionRule(comName, req.Command)
		}
		if rule != nil && ok {
			rules = append(rules, rule)
		}
	}
	for _, rule := range rules {
		if err := rule.Exec(&req, res); err != nil {
			return nil, err
		}
	}
	return res.Data, res.Error
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
	URI    string
	Code   string
	Domain string
}
