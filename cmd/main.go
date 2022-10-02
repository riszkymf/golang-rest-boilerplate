package main

import (
	// "fmt"

	"fmt"
	"log"
	"strings"

	"net/http"
	"os"

	"database/sql"

	"github.com/emicklei/go-restful/v3"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	route "github.com/riszkymf/golang-rest-boilerplate/internal"
	handler "github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	src "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type Env struct {
	DB_PATH    string
	APP_ENV    string
	APP_HOST   string
	APP_PORT   string
	WS_LOGGING string
	WS_AUTH    string
}

var env Env
var Connection *sql.DB

func init() {
	var err error
	if os.Getenv("APP_ENV") != "production" {
		godotenv.Load()
	}

	env.APP_ENV = src.GetEnv("APP_ENV", "staging")
	env.APP_HOST = src.GetEnv("APP_HOST", "0.0.0.0")
	env.APP_PORT = src.GetEnv("APP_PORT", "8080")
	env.DB_PATH = src.GetEnv("DB_PATH", "")
	env.WS_LOGGING = strings.ToUpper(src.GetEnv("WS_LOGGING", "TRUE"))
	env.WS_AUTH = strings.ToUpper(src.GetEnv("WS_AUTH", "TRUE"))

	if env.DB_PATH == "" {
		log.Fatal("DB Location is invalid")
	}

	Connection, err = sql.Open("sqlite3", env.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	handler.Connection = Connection

}

func main() {
	wsRConfig := route.RouteFilterConfig{
		WebServiceLogging: env.WS_LOGGING,
		Auth:              env.WS_AUTH,
	}
	ws := restful.NewContainer()
	ws = route.SetFilters(ws, wsRConfig)
	ws = route.SetRoutes(ws)
	appHost := fmt.Sprintf("%v:%v", env.APP_HOST, env.APP_PORT)
	log.Printf("Running App on %v", appHost)
	log.Fatal(http.ListenAndServe(appHost, ws))

}
