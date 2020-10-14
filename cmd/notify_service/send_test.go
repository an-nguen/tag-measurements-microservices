package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"testing"
)

func TestSendMail(t *testing.T) {
	var senderEmail = mail.Address{Address: "td-notify@dn-serv.ru"}
	var receiverEmail = mail.Address{Address: "an@dn-serv.ru"}
	var senderPass = "ifdxhgbrqmekzsti"
	var servername = "smtp.yandex.ru:465"
	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", senderEmail.Address, senderPass, host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	headers := make(map[string]string)
	headers["From"] = senderEmail.String()
	headers["To"] = receiverEmail.String()
	headers["Subject"] = "Test subj"
	body := "This is an example body.\n With two lines."
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(senderEmail.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(receiverEmail.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	_ = c.Quit()
}
