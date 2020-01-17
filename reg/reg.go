package reg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/context"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//code 到 RegObject 的缓存
var codeRegObjectMap map[string]*RegObject = make(map[string]*RegObject)

func GetRegistry() string {
	registry := configs.Default.App.Registry
	if registry == "" {
		registry = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
	}
	return registry
}
func localRegistry() string {
	return fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
}

func GetTokenContext(tokenCode string) (*context.Context, error) {
	//1、权限注册中心、2、应用注册中心，3、本地
	authAddr := ""
	if sers, err := FindByCode(configs.Default.Auth.Address); sers != nil && len(sers.Addrs[0]) > 0 {
		authAddr = sers.Addrs[0]
	} else {
		glog.Error(err)
	}
	if authAddr == "" {
		authAddr = configs.Default.App.Registry
	}
	if authAddr == "" {
		authAddr = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(authAddr)
	remoteUrl.Path = fmt.Sprintf("/api/oauth/token/%s", tokenCode)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var resBodyObj struct {
		Msg  string `json:"msg"`
		Data struct {
			AccessToken string `json:"access_token"`
			Type        string `json:"type"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return nil, err
	}
	token := &context.Context{}
	token, _ = token.FromTokenString(fmt.Sprintf("%s %s", resBodyObj.Data.Type, resBodyObj.Data.AccessToken))
	return token, nil
}
func RegisterDefault() error {
	addrs := make([]string, 0)
	if host := configs.Default.App.Address; host != "" {
		addrs = append(addrs, host)
	} else {
		ips := utils.GetIpAddrs()
		for _, item := range ips {
			addrs = append(addrs, fmt.Sprintf("http://%s:%s", item, configs.Default.App.Port))
		}
	}
	return Register(RegObject{
		Code:    configs.Default.App.Code,
		Name:    configs.Default.App.Name,
		Addrs:   addrs,
		Configs: configs.Default,
	})
}
func Register(item RegObject) error {
	client := &http.Client{}
	client.Timeout = 3 * time.Second
	postBody, err := json.Marshal(item)
	if err != nil {
		glog.Error(err)
		return err
	}
	regHost := GetRegistry()
	remoteUrl, _ := url.Parse(regHost)
	remoteUrl.Path = "/api/regs/register"
	req, err := http.NewRequest("POST", remoteUrl.String(), bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		glog.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return err
	}
	var resBodyObj struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return err
	}
	glog.Error("成功注册：", glog.Any("Addrs", item.Addrs), glog.Any("RegHost", regHost))
	return nil
}
func DoHttpRequest(serverName, method, path string, body io.Reader) ([]byte, error) {
	regs, err := FindByCode(serverName)
	if err != nil {
		return nil, err
	}
	if regs == nil || len(regs.Addrs) == 0 {
		return nil, glog.Error("找不到服务,", glog.String("serverName", serverName))
	}
	serverUrl := regs.Addrs[0]
	client := &http.Client{}
	remoteUrl, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	remoteUrl.Path = path
	req, err := http.NewRequest(method, remoteUrl.String(), body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		var resBodyObj struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
			return nil, err
		}
		return nil, glog.Error(resBodyObj.Msg)
	}
	return resBody, nil
}

func GetServerAddr(code string) (string, error) {
	if d, err := FindByCode(code); err != nil {
		return "", err
	} else if d != nil {
		return d.Addrs[0], nil
	}
	return "", nil
}

/**
通过编码找到注册对象
*/
func FindByCode(code string) (*RegObject, error) {
	if code == "" {
		return nil, nil
	}
	//优先从缓存里取
	ck := fmt.Sprintf("%s", strings.ToLower(code))
	if cv, ok := codeRegObjectMap[ck]; ok && cv != nil {
		return cv, nil
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(GetRegistry())
	remoteUrl.Path = fmt.Sprintf("/api/regs/%s", code)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		glog.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var resBodyObj struct {
		Msg  string     `json:"msg"`
		Data *RegObject `json:"data"`
	}
	glog.Error(string(resBody))
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		glog.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		glog.Error(resBodyObj.Msg)
		return nil, err
	}
	//设置缓存
	codeRegObjectMap[ck] = resBodyObj.Data

	return resBodyObj.Data, nil
}
