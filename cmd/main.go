package main

import (
	// "fmt"

	"fmt"
	"log"

	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	route "github.com/riszkymf/golang-rest-boilerplate/internal"
	handler "github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	src "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type Env struct {
	DB_PATH  string
	APP_ENV  string
	APP_HOST string
	APP_PORT string
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
	route.SetRoutes()
	appHost := fmt.Sprintf("%v:%v", env.APP_HOST, env.APP_PORT)
	log.Printf("Running App on %v", appHost)
	log.Fatal(http.ListenAndServe(appHost, nil))

}
