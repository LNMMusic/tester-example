package main

import (
	"fmt"
	"os"

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
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PSWD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DB_ADDR"),
			DBName:               os.Getenv("DB_NAME"),
		},
		Addr: os.Getenv("SERVER_ADDR"),
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