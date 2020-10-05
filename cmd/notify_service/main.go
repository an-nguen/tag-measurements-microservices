package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"Thermo-WH/internal/notify_service/structures"
	"Thermo-WH/pkg/datasource"
	"Thermo-WH/pkg/models"
)

var FILENAME = "/configs/config_notify.json"

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	var senderEmail = mail.Address{Address: "td-notify@dn-serv.ru"}
	var senderPass = "ifdxhgbrqmekzsti"
	var serverName = "smtp.yandex.ru:465"
	host, _, _ := net.SplitHostPort(serverName)
	headers := make(map[string]string)
	headers["From"] = senderEmail.String()
	headers["Subject"] = "Температурный режим WST-датчика"
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	auth := smtp.PlainAuth("", senderEmail.Address, senderPass, host)

	config := structures.ReadAppConfig(FILENAME)
	db := datasource.InitDatabaseConnection(config.Host, config.Port, config.User, config.Password, config.DbName)

	var groups []models.TemperatureZone
	for {
		// iterate groups
		db.Preload("Tags").Find(&groups)

		inMinute := time.Now()
		inMinute = inMinute.Add(-60 * time.Second)

		for _, group := range groups {
			lowerTempLimit := group.LowerTempLimit
			higherTempLimit := group.HigherTempLimit
			// if both limits equals null -> next warehouse group
			if lowerTempLimit == 0 && higherTempLimit == 0 {
				continue
			}
			toEmail := strings.Split(group.NotifyEmails, ",")

			// iterate tag managers
			go func(group models.TemperatureZone, db *gorm.DB) {
				for _, tag := range group.Tags {
					log.Printf("[LOG]\tНачинаем проверять тег %s\n", tag.Name)
					var rows []models.Measurement
					log.Printf("inMinute = %s", inMinute.String())
					db.Where("tag_uuid = ? and date > ? and (temperature < ? or temperature > ?)",
						tag.UUID, inMinute, lowerTempLimit, higherTempLimit).Find(&rows)

					// if rows is not null
					if len(rows) > 0 && len(toEmail) > 0 {
						log.Printf("[LOG]\tНайдены превышения лимитов... (%s)\n", tag.Name)
						for _, email := range toEmail {
							// prepare mail message
							receiverEmail := mail.Address{Address: email}
							headers["To"] = receiverEmail.String()
							message := ""

							body := ""
							var alert string
							for _, d := range rows {
								if d.Temperature > higherTempLimit {
									alert = fmt.Sprintf("[%s] Название тега - %s. Температура выше нормы (норма = %f C) - %f C",
										d.Date.Format("02.01.2006 15:04"), tag.Name, higherTempLimit, d.Temperature)
								} else if d.Temperature < lowerTempLimit {
									alert = fmt.Sprintf("[%s] Название тега - %s. Температура ниже нормы (норма = %f C) - %f C",
										d.Date.Format("02.01.2006 15:04"), tag.Name, lowerTempLimit, d.Temperature)
								}
								body += alert + "\n\n"
							}
							fmt.Printf("[LOG]\tОтправляем письмо получателю %s\n", email)

							for k, v := range headers {
								message += fmt.Sprintf("%s: %s\r\n", k, v)
							}
							message += "\r\n" + body
							// send mail
							conn, err := tls.Dial("tcp", serverName, tlsconfig)
							if err != nil {
								log.Print(err)
								continue
							}

							c, err := smtp.NewClient(conn, host)
							if err != nil {
								log.Print(err)
								continue
							}

							// Auth
							if err = c.Auth(auth); err != nil {
								log.Print(err)
								continue
							}

							// To && From
							if err = c.Mail(senderEmail.Address); err != nil {
								log.Print(err)
								continue
							}

							if err = c.Rcpt(receiverEmail.Address); err != nil {
								log.Print(err)
								continue
							}

							// Data
							w, err := c.Data()
							if err != nil {
								log.Print(err)
								continue
							}

							_, err = w.Write([]byte(message))
							if err != nil {
								log.Print(err)
								continue
							}

							err = w.Close()
							if err != nil {
								log.Print(err)
								continue
							}

							_ = c.Quit()
						}
					} else {
						continue
					}
				}
			}(group, db)

		}
		time.Sleep(time.Second * 60)
	}
}
