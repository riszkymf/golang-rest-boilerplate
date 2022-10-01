package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	"github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type Env struct {
	DB_PATH  string
	APP_ENV  string
	APP_HOST string
	APP_PORT string
}

var env Env

type FuncTest func(*testing.T, *sql.DB)

func TestDBFunctionality(t *testing.T) {
	var err error
	var Connection *sql.DB
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

	dataHolder := map[string][]int{
		"author": {6},
	}

	Connection, err = sql.Open("sqlite3", env.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	handler.Connection = Connection

	t.Run("Setup Data", func(t *testing.T) {
		var err error

		inputAuthor := map[string]any{
			"id":   6,
			"name": "Herman Melville",
		}
		inputAuthor["id"], err = handler.InsertData("author", inputAuthor)
		if err != nil {
			t.Fatalf(`Error: %v`, err.Error())
		}
		inputBook := []map[string]any{
			{
				"title":     "moby dick",
				"stock":     14,
				"author_id": inputAuthor["id"],
			},
			{
				"title":     "bartleby, the scrivener",
				"stock":     5,
				"author_id": inputAuthor["id"],
			},
			{
				"title":     "benito cereno",
				"stock":     5,
				"author_id": inputAuthor["id"],
			},
			{
				"title":     "isle of the cross",
				"stock":     7,
				"author_id": inputAuthor["id"],
			},
		}

		bookResult, err := handler.InsertMultipleData("books", inputBook)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		fmt.Println("Insert Result : ", bookResult)
		dataHolder["books"] = bookResult

		inputMembers := []map[string]any{
			{
				"email":     "janice@email.com",
				"firstname": "janice",
				"lastname":  "sopranos",
				"address":   "newark",
			},
			{
				"email":     "hidayat@email.com",
				"firstname": "hidayat",
				"lastname":  "soekamto",
				"address":   "jogja",
			},
			{
				"email":     "fgibbs@email.com",
				"firstname": "freddie",
				"lastname":  "gibbs",
				"address":   "gary",
			},
		}
		result, err := handler.InsertMultipleData("members", inputMembers)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		dataHolder["members"] = result
		fmt.Println("Insert Member Result : ", result)
	})

	booksData := []map[string]any{}

	t.Run("Query", func(t *testing.T) {

		testBuild := handler.FilterQuery{
			And: map[string][]handler.FieldFilter{
				"author_name": {
					{
						Operator:  "like",
						Value:     "Herman%",
						ValueType: "string",
					},
					{
						Operator:  "like",
						Value:     "%Melville",
						ValueType: "string",
					},
				},
				"stock": {
					{
						Operator:  "lt",
						Value:     "10",
						ValueType: "int",
					},
				},
			},
			Or: map[string][]handler.FieldFilter{
				"title": {
					{
						Operator:  "like",
						Value:     "bartleby%",
						ValueType: "string",
					},
					{
						Operator:  "like",
						Value:     "%cross",
						ValueType: "string",
					},
				},
			},
		}

		result, err := handler.GetRowByFilter("v_books", testBuild)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		fmt.Println("Insert Result : ", result)
		booksData = result

	})

	t.Run("Update and transaction", func(t *testing.T) {
		members, err := handler.GetRowsAll("members")
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		fmt.Println(members)
		fmt.Println(booksData)
		insertRecordData := map[string]any{
			"book_id":     booksData[0]["book_id"],
			"member_id":   members[3]["id"],
			"rent_date":   time.Now().Local().Format("2006-01-02"),
			"due_date":    time.Now().AddDate(0, 0, 14).Local().Format("2006-01-02"),
			"rent_status": "rented",
		}
		fmt.Println(insertRecordData)
		res, err := handler.InsertData("records", insertRecordData)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		fmt.Println("Record inserted with id ", res)

		bookStock, ok := booksData[0]["stock"].(int)
		if !ok {
			t.Fatalf(`Error: type is not number`)
		}
		updateBookData := map[string]any{
			"stock": bookStock - 1,
		}
		bookId := booksData[0]["book_id"].(int)
		err = handler.UpdateData("books", updateBookData, bookId)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		err := handler.DeleteData("author", 6)
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
		err = handler.DeleteMultipleData("members", dataHolder["members"])
		if err != nil {
			t.Fatalf(`Error: %v`, err)
		}
	})

}
