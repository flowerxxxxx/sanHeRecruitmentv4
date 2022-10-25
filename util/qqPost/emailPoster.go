package qqPost

import (
	"gopkg.in/gomail.v2"
	"sanHeRecruitment/config"
)

func SendEmail(reciever, code string) (err error) {
	m := gomail.NewMessage()
	//发送人
	m.SetHeader("From", "2633122565@qq.com")
	//接收人
	m.SetHeader("To", reciever)
	//抄送人
	m.SetAddressHeader("错误娘", "sanHeRec", "三河招聘平台")
	//主题
	m.SetHeader("Subject", "Code")
	//内容
	body := "<h3>您的验证码为：" + code + "，该验证码在十分钟内有效，请注意时间哦！</h3>"
	m.SetBody("text/html", body)
	m.Attach(config.ErrorLogAddr)
	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, "2633122565@qq.com", "raaepljblyzydiaa")
	// 发送邮件
	err = d.DialAndSend(m)
	return err
}
