package smtp

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"os"
)

// Mail ...
type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

// SMTPServer ...
type SMTPServer struct {
	Host      string
	Port      string
	Password  string
	Login     string
	TLSConfig *tls.Config
}

var (
	sender = SMTPServer{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Login:    os.Getenv("SMTP_LOGIN"),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.sendgrid.com"}}

	mail = Mail{
		Sender:  "butler@jenkins.io",
		To:      []string{"me@olblak.com"},
		Subject: "Jenkins Election Participation List",
		Body:    "Hello my body"}
)

func sendEmail(body string) {

	auth := smtp.PlainAuth("", sender.Login, sender.Password, sender.Host)

	conn, err := tls.Dial("tcp", sender.Host+":"+sender.Port, sender.TLSConfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, sender.Host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		log.Panic(err)
	}
	receivers := mail.To
	for _, k := range receivers {
		log.Println("sending to: ", k)
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	mail.Body = body

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(mail.Body))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
