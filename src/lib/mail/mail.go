package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

func encodeRFC2047(String string) string {
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

/**
 * 邮件发送
 * @param toMail  		邮件接收者账号
 * @param Subject		邮件主题
 * @param fromName		邮件发送着名称（便于接收者直观查看邮件来源)
 * @param mailtype      邮件内容类型： html | plain
 * @param body			邮件内容
 *
 *
 */
func SendMail(toMail, Subject, fromName, mailtype, body string) error {

	smtpServer := "smtpdm.aliyun.com"
	auth := smtp.PlainAuth(
		"",
		"customerservice@mail.starrymobi.com",
		"O0EfM3fWmiqjUyRg",
		"smtpdm.aliyun.com",
	)
	from := mail.Address{fromName, "customerservice@mail.starrymobi.com"}
	to := mail.Address{"Receiver", toMail}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = Subject
	//header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/" + mailtype + "; charset=\"utf-8\""
	//header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	fmt.Println("start send mail to : " + toMail)
	//err := smtp.SendMail(
	//	smtpServer+":465",
	//	auth,
	//	from.Address,
	//	[]string{to.Address},
	//	[]byte(message),
	//)

	err := SendMailUsingTLS(
		smtpServer+":465",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	fmt.Println("finished sen mail.")

	return err
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
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
