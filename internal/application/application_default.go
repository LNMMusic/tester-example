package application

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/LNMMusic/tester-example/internal/task/handler"
	"github.com/LNMMusic/tester-example/internal/task/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

var (
	// ErrApplicationRun is returned when the application fails to run
	ErrApplicationRun = errors.New("application: failed to run")
	// ErrApplicationSetUp is returned when the application fails to set up
	ErrApplicationSetUp = errors.New("application: failed to set up")
	// ErrApplicationTearDown is returned when the application fails to tear down
	ErrApplicationTearDown = errors.New("application: failed to tear down")
)

// ConfigApplicationDefault is the configuration of ApplicationDefault
type ConfigApplicationDefault struct {
	// database connection
	DbCfg *mysql.Config
	// Addr is the address where the server will listen
	Addr string
}

// NewApplicationDefault returns a new ApplicationDefault
func NewApplicationDefault(cfg *Config) *ApplicationDefault {
	// default config
	defaultConfig := &Config{
		Db: (*mysql.Config)(nil),
		Addr: ":8080",
	}
	if cfg != nil {
		if cfg.Db != nil {
			defaultConfig.Db = cfg.Db
		}
		if cfg.Addr != "" {
			defaultConfig.Addr = cfg.Addr
		}
	}

	return &ApplicationDefault{
		Config: defaultConfig,
		rt: chi.NewRouter(),
		db: (*sql.DB)(nil),
	}
}

// Config is the configuration of ApplicationDefault
type Config struct {
	// database connection
	Db *mysql.Config
	// Addr is the address where the server will listen
	Addr string
}

// ApplicationDefault is the default implementation of Application interface
type ApplicationDefault struct {
	// configuration of the application
	*Config

	// instances of the application
	// - db is the database connection
	db *sql.DB
	// - rt is the router using chi
	rt *chi.Mux
}

// TearDown tears down the application
// - close resources in reverse order
func (a *ApplicationDefault) TearDown() (err error) {
	// close database connection
	if a.db != nil {
		err = a.db.Close()
		if err != nil {
			err = fmt.Errorf("%w. %v", ErrApplicationTearDown, err)
			return
		}
		return
	}
	return
}

// SetUp sets up the application
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - database: connection
	a.db, err = sql.Open("mysql", a.Config.Db.FormatDSN())
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrApplicationSetUp, err)
		return
	}
	// - database: ping
	err = a.db.Ping()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrApplicationSetUp, err)
		return
	}
	// - storage
	st := storage.NewTaskMySQL(a.db)
	// - handler
	hd := handler.NewTask(st)

	// router
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	// - routes / endpoints
	a.rt.Get("/tasks/{id}", hd.GetById())
	a.rt.Post("/tasks", hd.Create())

	return
}
// Run runs the application
func (a *ApplicationDefault) Run() (err error) {
	err = http.ListenAndServe(a.Config.Addr, a.rt)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrApplicationRun, err)
		return
	}

	return
}