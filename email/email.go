package email

import (
	"crypto/tls"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes"

	"strings"
)

var emailConfig *EmailParam
var m *gomail.Message

func SendEmailNew(myEmail *EmailParam, targetAccount []string, content string, hermesOps *hermes.Hermes, subject string) {
	myToers := strings.Join(targetAccount, ",")
	myCCers := ""
	emailConfig = myEmail

	myEmail.Toers = myToers
	myEmail.CCers = myCCers

	h := hermesOps

	email := hermes.Email{
		Body: hermes.Body{
			Name: "尊敬的用户",
			Intros: []string{
				content,
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	InitEmail(myEmail)
	SendEmail(subject, emailBody)

}

func InitEmail(ep *EmailParam) {
	toers := []string{}

	m = gomail.NewMessage()

	if len(ep.Toers) == 0 {
		fmt.Println(len(ep.Toers))
		return
	}

	for _, tmp := range strings.Split(ep.Toers, ",") {
		toers = append(toers, strings.TrimSpace(tmp))
	}
	fmt.Println(len(ep.Toers))

	// 收件人可以有多个，故用此方式
	m.SetHeader("To", toers...)

	//抄送列表
	if len(ep.CCers) != 0 {
		for _, tmp := range strings.Split(ep.CCers, ",") {
			toers = append(toers, strings.TrimSpace(tmp))
		}
		m.SetHeader("Cc", toers...)
	}
	fmt.Println(len(ep.Toers))
	// 发件人
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", ep.FromEmail, "")
}

type EmailParam struct {
	// ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
	ServerHost string
	// ServerPort 邮箱服务器端口，如腾讯企业邮箱为465
	ServerPort int
	// FromEmail　发件人邮箱地址
	FromEmail string
	// FromPasswd 发件人邮箱密码（注意，这里是明文形式），TODO：如果设置成密文？
	FromPasswd string
	// Toers 接收者邮件，如有多个，则以英文逗号(“,”)隔开，不能为空
	Toers string
	// CCers 抄送者邮件，如有多个，则以英文逗号(“,”)隔开，可以为空
	CCers string
}

// SendEmail body支持html格式字符串
func SendEmail(subject, body string) {
	// 主题
	m.SetHeader("Subject", subject)

	// 正文
	m.SetBody("text/html", body)
	fmt.Println(m)
	d := gomail.NewDialer(emailConfig.ServerHost, emailConfig.ServerPort, emailConfig.FromEmail, emailConfig.FromPasswd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 发送
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
}
