package messageUtil

import (
	"sanHeRecruitment/util/timeUtil"
	"time"
)

// UpgradeApplySuccessTem 身份升级成功模板
func UpgradeApplySuccessTem(name, gender string, applyTime time.Time) (templateCall string) {
	timeStr := timeUtil.TimeFormatToStr(applyTime)
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"您于 " + timeStr + "发起的身份升级已审核通过，您需要通过'我的-切换身份'来切换到[企业端/服务机构端]以激活身份" +
		"。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

// UpgradeApplyErrTem 身份升级失败模板
func UpgradeApplyErrTem(name, gender, rejectReason string, applyTime time.Time) (templateCall string) {
	timeStr := timeUtil.TimeFormatToStr(applyTime)
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"您于 " + timeStr + "发起的身份升级已审核未通过，审核失败原因：" + rejectReason +
		"。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

// ResumeManageTem 简历审核 通过/未通过 模板
func ResumeManageTem(name, gender, ReqTitle string, applyTime time.Time, qualification int) (templateCall string) {
	timeStr := timeUtil.TimeFormatToStr(applyTime)
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	ifPass := "?"
	if qualification == 1 {
		ifPass = "通过"
	} else if qualification == 2 {
		ifPass = "未通过"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"您于 " + timeStr + "对'" + ReqTitle + "'投递的简历初审" + ifPass +
		"。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	if qualification == 1 {
		templateCall = templateCall + "请尽快与该公司联系，或通过小程序的在线聊天功能留下联系方式。"
	}
	return
}

// BossPubPassOrNotTem boss发布的招聘/需求审核 通过/未通过 模板
func BossPubPassOrNotTem(name, gender, ReqTitle, reason string, pubTime time.Time, qualification int) (templateCall string) {
	timeStr := timeUtil.TimeFormatToStr(pubTime)
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	ifPass := "?"
	if qualification == 1 {
		ifPass = "审核通过"
	} else if qualification == 2 {
		ifPass = "审核未通过" + "，未通过原因为[" + reason + "]"
	}

	templateCall = "亲爱的" + name + caller + "，" +
		"您于 " + timeStr + "发布的'" + ReqTitle + "'管理员" + ifPass +
		"。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

func BossBeDel(name, gender string) (templateCall string) {
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"您已被管理员重置为普通用户身份" +
		"。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

// AdminDeletePubInfo 管理员删除文章通知模板
func AdminDeletePubInfo(name, gender, ReqTitle, reason string, pubTime time.Time, qualification int) (templateCall string) {
	timeStr := timeUtil.TimeFormatToStr(pubTime)
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"您于 " + timeStr + "发布的'" + ReqTitle + "'已被管理员删除，删除原因为【" + reason +
		"】。处理时间：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

// InviteDeliveryTem 邀请投递消息推送模板
func InviteDeliveryTem(name, gender, companyName, bossName, ReqTitle string) (templateCall string) {
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" +
		"来自" + companyName + "的" + bossName + "对您发起了邀请投递，岗位名称为'" + ReqTitle +
		"',您可通过【我的-求职进展-邀请投递】对邀请进行查看" +
		"。邀请的发起时间为：" + time.Now().Format("2006-01-02 15:04:05")
	return
}

// DeliveryPushTem 简历投递模板
func DeliveryPushTem(name, gender, fromNickName, title string) (templateCall string) {
	caller := ""
	if gender == "男" {
		caller = "先生"
	} else {
		caller = "女士"
	}
	templateCall = "亲爱的" + name + caller + "，" + fromNickName + "对您所发布的'" + title +
		"'投递了简历"
	return
}
