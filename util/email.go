package util

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"net/textproto"
	"time"
)

// GenerateVerificationCode 生成验证码
func GenerateVerificationCode() string {
	// 生成 6 位数验证码
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06v", rand.Intn(1000000))
	fmt.Println("code:", code)
	return code

}

// SendVerificationCode 发送邮箱验证码
func SendVerificationCode(email1 string, code string) error {
	host := "smtp.qq.com" // qq邮箱smtp服务器地址
	//port := "25"  //非SSL协议端口
	port := "465" //SSl协议端口
	userName := "835118160@qq.com"
	password := "hiombizscucvbfae" // qq邮箱填授权码（非密码）

	e := &email.Email{
		To:      []string{email1},
		From:    userName,
		Subject: "Email Send Test",                                 // 发送内容标题
		Text:    []byte("Text Body is, of course, supported!"),     //信息内容
		HTML:    []byte("<h1>This is your code:" + code + "</h1>"), //信息内容 （html格式）
		Headers: textproto.MIMEHeader{},
	}
	// 非SSL协议端口25
	//err := e.Send(host+":"+port, smtp.PlainAuth("", userName, password, host))
	//if err != nil {
	//	panic(err)
	//}
	// 使用SSl协议端口 465/587
	return e.SendWithTLS(host+":"+port, smtp.PlainAuth("", userName, password, host), &tls.Config{InsecureSkipVerify: true, ServerName: host})
}
