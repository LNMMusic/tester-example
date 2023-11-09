package main

import (
	"fmt"

	"github.com/LNMMusic/tester-example/internal/application"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// env
	// ...

	// application
	// - config
	cfg := &application.Config{
		Db: &mysql.Config{
			User:                 "root",
			Passwd:               "",
			Net:                  "tcp",
			Addr:                 "localhost:3306",
			DBName:               "tester_example_tasks_db",
		},
		Addr: "localhost:8080",
	}
	// - new
	app := application.NewApplicationDefault(cfg)
	// - tear down
	defer app.TearDown()
	// - set up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}