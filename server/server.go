package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/danilotadeu/star_wars/api"
	"github.com/danilotadeu/star_wars/app"
	"github.com/danilotadeu/star_wars/store"
)

// Server is a interface to define contract to server up
type Server interface {
	Start()
	ConnectDatabase() *sql.DB
}

type server struct {
	App   *app.Container
	Store *store.Container
	Db    *sql.DB
}

// New is instance the server
func New() Server {
	return &server{}
}

func (e *server) Start() {
	e.Db = e.ConnectDatabase()
	e.Store = store.Register(e.Db, os.Getenv("URL_STARWARS_API"))
	e.App = app.Register(e.Store)
	api.Register(e.App, os.Getenv("PORT"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = e.Db.Close()
	}()
}

func (e *server) ConnectDatabase() *sql.DB {
	connectionMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", connectionMysql)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		log.Println("error db.Ping(): ", err.Error())
		panic(err)
	}

	return db
}
