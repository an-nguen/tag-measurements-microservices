package datasource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
)

/**
 *	Function: initDatabaseConnection
 *  --------------------------------
 *  Create database connection by passing AppConfig
 *
 *	returns: sql.DB struct
 **/
func InitDatabaseConnection(host string, port string, user string, pass string, dbname string) *gorm.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbname)
	db, err := gorm.Open("postgres", connStr)
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
