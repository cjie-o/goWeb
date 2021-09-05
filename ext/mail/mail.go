package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

type Mail struct {
	User string
	Pwd  string
	From string
	To   string
	Sub  string
	Body string
}

func NewMail() *Mail {
	return &Mail{}
}
func (M *Mail) Init(User string, Pwd string, From string) *Mail {
	M.User = User
	M.Pwd = Pwd
	M.From = From
	return M
}
func (M *Mail) Set(To string, Sub string, Body string) *Mail {
	M.To = To
	M.Sub = Sub
	M.Body = Body
	return M
}
func (M *Mail) Send() error {
	host := "smtp.exmail.qq.com"
	port := 465
	email := M.User
	pwd := M.Pwd    // 这里填你的授权码
	toEmail := M.To // 目标地址

	header := make(map[string]string)

	header["From"] = M.From + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = M.Sub
	header["Content-Type"] = "text/html;chartset=UTF-8"

	body := M.Body

	message := ""

	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		email,
		pwd,
		host,
	)

	err := M.sendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)
	return err
	// if err != nil {
	// 	panic(err)
	// }

}

//return a smtp client
func (M *Mail) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func (M *Mail) SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func (M *Mail) sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := M.dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}