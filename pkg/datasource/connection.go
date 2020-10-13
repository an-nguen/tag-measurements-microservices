package datasource

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**
 *	Function: initDatabaseConnection
 *  --------------------------------
 *  Create database connection by passing AppConfig
 *
 *	returns: sql.DB struct
 **/
func InitDatabaseConnection(host string, port string, user string, pass string, dbname string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

/**
 *	Function: initDatabaseConnection
 *  --------------------------------
 *  Create database connection by passing AppConfig
 *
 *	returns: sql.DB struct
 **/
func InitAmqpConn(uri string) *amqp.Channel {
	conn, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return channel
}
