package wechatPubAcc

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	APPID     = "wx5f80a46aee1c9bdd"
	APPSECRET = "c1acfdce44c49c34b77ad88c5b690d16"
	//SentTemplateID = "Al8FCd4p2gIFx1KrrlTJprM_twK6Fzn7CItzrHHXgvU" //每日一句的模板ID，替换成自己的
	ConversationMessageTemplateID = "sUor3v4Ve_0T3QnuiYSjUkc6wB5oqdW7L4vuOjzvJ2k" //消息通知模板id
	DeliveryResumeTemplateID      = "xmmR-qPCuwX28cxos4acpXdRzjSVceJpEnrnBMrg7_0"
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
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)

	//url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", config.WechatAppid, config.WechatSecret)
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
