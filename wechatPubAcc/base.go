package wechatPubAcc

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"strings"
	"time"
)

var (
	APPID     = config.WechatPublicAppid
	APPSECRET = config.WechatPublicSecret
	//SentTemplateID = "Al8FCd4p2gIFx1KrrlTJprM_twK6Fzn7CItzrHHXgvU" //每日一句的模板ID，替换成自己的
	ConversationMessageTemplateID = config.WechatConversationMessageTemplateID //消息通知模板id
	DeliveryResumeTemplateID      = config.WechatDeliveryResumeTemplateID
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type sentence struct {
	Content     string `json:"content"`
	Note        string `json:"note"`
	Translation string `json:"translation"`
}

// 获取微信accesstoken
func getaccesstoken() string {
	//url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
	//url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", config.WechatAppid, config.WechatSecret)
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return ""
	}
	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}
	ATSaveTime := time.Duration(token.ExpiresIn) * time.Second
	if errReSet := dao.Redis.Set("wechat_public_access_token", token.AccessToken, ATSaveTime).Err(); errReSet != nil {
		log.Println("getaccesstoken set failed,err:", errReSet)
		return ""
	}

	fmt.Println(token)
	fmt.Println("微信token解析json", token.AccessToken, token.ExpiresIn)

	return token.AccessToken
}

// 获取当前微信公众号关注者列表
func getflist(access_token string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	flist := gjson.Get(string(body), "data.openid").Array()
	//fmt.Println(flist)
	return flist
}

// 发送模板消息 //携带访问url页面
func templatepostUrl(access_token string, reqdata string, fxurl string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"url\":\"" + fxurl + "\", \"data\": " + reqdata + "}"

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// 发送模板消息 //不携带访问url页面
func templatepost(access_token string, reqdata string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token
	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\",\"data\": " + reqdata + "}"
	//fmt.Println(reqbody)
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		log.Println("httpPostErr:", err)
		return
	}
	//fmt.Println("templatepost err:", err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("readErr:", body, err)
		return
	}
	//fmt.Println("消息发送成功")
}

//func templatepost(access_token string, reqdata string, templateid string, openid string) {
//	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token
//	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\",\"data\": " + reqdata + "}"
//	//fmt.Println(reqbody)
//	resp, err := http.Post(url,
//		"application/x-www-form-urlencoded",
//		strings.NewReader(string(reqbody)))
//	if err != nil {
//		log.Println("httpPostErr:", err)
//		return
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Println("readErr:", body, err)
//		return
//	}
//	//fmt.Println("消息发送成功")
//}
