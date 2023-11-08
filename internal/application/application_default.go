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
func NewApplicationDefault(cfg *ConfigApplicationDefault) *ApplicationDefault {
	// default config
	defaultRt := chi.NewRouter()
	defaultDb := (*sql.DB)(nil)
	defaultCfg := &ConfigApplicationDefault{
		Addr: ":8080",
	}
	if cfg != nil {
		if cfg.DbCfg != nil {
			defaultCfg.DbCfg = cfg.DbCfg
		}
		if cfg.Addr != "" {
			defaultCfg.Addr = cfg.Addr
		}
	}

	return &ApplicationDefault{
		dbCfg: defaultCfg.DbCfg,
		db: defaultDb,
		rt: defaultRt,
		addr: defaultCfg.Addr,
	}
}

// ApplicationDefault is the default implementation of Application interface
type ApplicationDefault struct {
	// database connection
	dbCfg *mysql.Config
	db *sql.DB
	// rt is the router using chi
	rt *chi.Mux
	// addr is the address where the server will listen
	addr string
}

// TearDown tears down the application
// - close resources in reverse order
func (a *ApplicationDefault) TearDown() (err error) {
	// close database connection
	defer func() {
		if a.db != nil {
			err = a.db.Close()
			if err != nil {
				err = fmt.Errorf("%w. %v", ErrApplicationTearDown, err)
				return
			}
		}
	}()

	return
}

// SetUp sets up the application
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - database connection
	a.db, err = sql.Open("mysql", a.dbCfg.FormatDSN())
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
	err = http.ListenAndServe(a.addr, a.rt)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrApplicationRun, err)
		return
	}

	return
}