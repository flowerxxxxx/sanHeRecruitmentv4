package wechatPubAcc

import (
	"sanHeRecruitment/models/mysqlModel"
)

// ConversationMessagePush 向退出小程序并关注公众号的用户推送会话消息
func ConversationMessagePush(openid, fromUser, content string) {
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
	reqdata := "{\"fromUser\":{\"value\":\"" + "来自" + fromUser + "发送的消息" + "\", \"color\":\"#0000CD\"}, \"message\":{\"value\":\"" + "消息内容:" + content + "\"}, \"intention\":{\"value\":\"" + "请前往小程序查看消息" + "\"}}"
	templatepost(access_token, reqdata, ConversationMessageTemplateID, pubId)
}

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
