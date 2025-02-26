package auth

import (
	"gopkg.in/gomail.v2"
	//"log"
	//"fmt"
)

func Mailing(subject string, body *string, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "Ecommerce@gmail.com")
	m.SetHeader("To", to)
	m.SetAddressHeader("Cc", "phychemma4@gmail.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", *body)
	//	m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 465, "phychemma4@gmail.com", "zzku acdb tsvs svgp")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
