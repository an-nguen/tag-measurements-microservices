package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"tag-measurements-microservices/internal/tgbot_service/structures"
	"tag-measurements-microservices/pkg/datasource"
	"tag-measurements-microservices/pkg/models"
	"time"
)

var FILENAME = "/configs/config_tgbot.json"

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	config := structures.ReadAppConfig(FILENAME)
	db := datasource.InitDatabaseConnection(config.Host, config.Port, config.User, config.Password, config.DbName)
	var groupId = config.GroupId

	bot, err := tgbotapi.NewBotAPI(config.TGBotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	go func(db *gorm.DB) {
		var groups []models.TemperatureZone
		for {
			db.Preload("Tags").Find(&groups)
			for _, group := range groups {
				lowerTempLimit := group.LowerTempLimit
				higherTempLimit := group.HigherTempLimit
				// if both limits equals null -> next warehouse group
				if lowerTempLimit == 0 && higherTempLimit == 0 {
					continue
				}
				inMinute := time.Now()
				inMinute = inMinute.Add(-45 * time.Second)
				// iterate tag managers
				go func(group models.TemperatureZone, db *gorm.DB) {
					for _, tag := range group.Tags {
						log.Printf("[LOG]\tНачинаем проверять тег %s\n", tag.Name)
						var rows []models.Measurement
						log.Printf("[LOG]\tinMinute = %s", inMinute.String())
						db.Where("tag_uuid = ? and date > ? and (temperature < ? or temperature > ?)",
							tag.UUID, inMinute, lowerTempLimit, higherTempLimit).Find(&rows)

						// if rows is not null
						if len(rows) > 0 {
							log.Printf("[LOG]\tНайдены превышения температурных норм... (%s)\n", tag.Name)
							// prepare message
							var alert string
							for _, d := range rows {
								if d.Temperature > higherTempLimit {
									alert = fmt.Sprintf("[%s] Тег '%s' - температура выше нормы (норма = %f) < %f",
										d.Date.Format("02.01.2006 15:04"), tag.Name, higherTempLimit, d.Temperature)
								} else if d.Temperature < lowerTempLimit {
									alert = fmt.Sprintf("[%s] Тег '%s' - температура ниже нормы (норма = %f) < %f",
										d.Date.Format("02.01.2006 15:04"), tag.Name, lowerTempLimit, d.Temperature)
								}
							}
							fmt.Printf("[LOG]\tОтправляем сообщение в телеграм группу с ид %d.\n", groupId)
							msg := tgbotapi.NewMessage(groupId, alert)
							_, _ = bot.Send(msg)
						}
					}
				}(group, db)
			}
			time.Sleep(time.Second * 45)
		}
	}(db)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(groupId, "")
			switch update.Message.Command() {
			case "start_bot2":
				msg.Text = "Это второй бот для отправки уведомлении о температурном превышении/понижение нормы."
			case "sayhi":
				msg.Text = "Привет :)"
			}
			_, _ = bot.Send(msg)
		}
	}
}
