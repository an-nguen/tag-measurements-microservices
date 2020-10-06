package datasource

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	_ "github.com/lib/pq"

	utils "tag-measurements-microservices/internal/fetch_service/structures"
)

var wg sync.WaitGroup
var db *sql.DB

func routine(name string) {
	var i int
	err := db.QueryRow("INSERT INTO public.warehouse_group (name, description) VALUES ($1, 'TEST')", name).Scan(&i)
	if err != nil {
		panic(err)
	}
	fmt.Println(i)

	wg.Add(-1)
}

func TestMain(m *testing.M) {
	// Read config
	var appConfig = utils.ReadAppConfig("/configs/config_fetch.json")

	db = InitDatabaseConnection(appConfig.Host, appConfig.Port, appConfig.User, appConfig.Password, appConfig.DbName)
	wg.Add(3)

	go routine("TEST1")
	go routine("TEST2")
	go routine("TEST3")

	wg.Wait()
}
