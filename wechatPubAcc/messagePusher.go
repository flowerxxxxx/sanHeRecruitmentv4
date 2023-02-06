package wechatPubAcc

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"time"
)

// ConversationMessagePush 向退出小程序并关注公众号的用户推送会话消息
func ConversationMessagePush(openid, fromUser, content string) {
	//获取公众号的token
	access_token := dao.Redis.Get("wechat_public_access_token").Val()
	if access_token == "" {
		access_token = getaccesstoken()
		if access_token == "" {
			//不推送或推送失败
			return
		}
	}
	//获取被推送用户的id
	//pubId := mysqlModel.OpenToPub(openid)
	//if pubId == "" {
	//	//不推送或推送失败
	//	return
	//}
	reqdata := "{\"first\":{\"value\":\"" + "您好，您收到一个人选意向沟通请查看。" + "\", \"color\":\"#0000CD\"}," +
		" \"keyword1\":{\"value\":\"" + fromUser + "\"}," +
		" \"keyword2\":{\"value\":\"" + "待沟通" + "\"}," +
		" \"keyword3\":{\"value\":\"" + "进行中" + "\" }," +
		" \"keyword3\":{\"value\":\"" + time.Now().Format("2006-01-02 15:04:05") + "\"}," +
		" \"remark\" : {\"value\":\"" + "沟通内容：" + content + "\"} }"
	templatepost(access_token, reqdata, ConversationMessageTemplateID, openid)
}

//func ConversationMessagePush(openid, fromUser, content string) {
//	//获取公众号的token
//	access_token := getaccesstoken()
//	if access_token == "" {
//		//不推送或推送失败
//		return
//	}
//	//获取被推送用户的id
//	//pubId := mysqlModel.OpenToPub(openid)
//	//if pubId == "" {
//	//	//不推送或推送失败
//	//	return
//	//}
//	reqdata := "{\"name1\":{\"value\":\"" + fromUser + "\"}, \"thing3\":{\"value\":\"" + content + "\"}, \"thing4\":{\"value\":\"" + "消息推送" + "\"},\"time5\":{\"value\":\"" + time.Now().Format("2006-01-02 15:04:05") + "\"}}"
//	templatepost(access_token, reqdata, "NsmDzZKmyMsLlsQDO9X9c62S5vQsFt66rS8NI1EmQcA", openid)
//}

// DeliveryResumeMessagePush 向退出小程序并关注公众号的招聘者推送招聘简历投递消息
func DeliveryResumeMessagePush(openid, fromUser, content string) {
	//获取公众号的token
	access_token := getaccesstoken()
	if access_token == "" {
		//不推送或推送失败
		return
	}
	//获取被推送用户的id
	pubId := mysqlModel.OpenToPub(openid)
	if pubId == "" {
		//不推送或推送失败
		return
	}
	reqdata := "{\"fromUser\":{\"value\":\"" + fromUser + "投递了简历" + "\", \"color\":\"#0000CD\"}, \"message\":{\"value\":\"" + fromUser + "向您所发布的'" + content + "'招聘需求投递了简历" + "\"}, \"intention\":{\"value\":\"" + "请前往小程序查看消息" + "\"}}"
	templatepost(access_token, reqdata, DeliveryResumeTemplateID, pubId)
}

//发送每日一句
//func everydaysen() {
//
//	req, fxurl := getsen()
//	if req.Content == "" {
//		return
//	}
//	access_token := getaccesstoken()
//	if access_token == "" {
//		return
//	}
//
//	flist := getflist(access_token)
//	if flist == nil {
//		return
//	}
//
//	reqdata := "{\"content\":{\"value\":\"" + req.Content + "\", \"color\":\"#0000CD\"}, \"note\":{\"value\":\"" + req.Note + "\"}, \"translation\":{\"value\":\"" + req.Translation + "\"}}"
//	//遍历关注列表，全员发送
//	for _, v := range flist {
//		templatepostUrl(access_token, reqdata, fxurl, SentTemplateID, v.Str)
//	}
//}
