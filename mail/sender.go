package mail

import "gopkg.in/gomail.v2"

const (
	smtpAddress = "smtp.gmail.com"
	smtpPort    = 587
)

type EmailSender interface {
	SendEmail(
		to []string,
		subject string,
		body string,
	) error
}

type GmailSender struct {
	Name     string
	Address  string
	Password string
}

func NewGmailSender(name, address, password string) EmailSender {
	return &GmailSender{
		Name:     name,
		Address:  address,
		Password: password,
	}
}

func (sender *GmailSender) SendEmail(to []string, subject string, body string) error {
	return nil
	e := gomail.NewMessage()

	e.SetHeader("From", sender.Address)
	e.SetHeader("To", to...)
	e.SetBody("text/html", body)

	d := gomail.NewDialer(smtpAddress, smtpPort, sender.Address, sender.Password)

	if err := d.DialAndSend(e); err != nil {
		return err
	}
	return nil
}
