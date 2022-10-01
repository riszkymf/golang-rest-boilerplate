package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"strconv"

	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

var Connection *sql.DB

type FilterQuery struct {
	And map[string][]FieldFilter

	Or map[string][]FieldFilter
}

type FieldFilter struct {
	Operator  string
	Value     string
	ValueType string
}

func buildFilter(filterObj FilterQuery) (string, error) {
	andFilters := []string{}
	orFilters := []string{}
	var andFilter, orFilter string
	if filterObj.And != nil {
		for k, v := range filterObj.And {
			filter, err := buildFilterString(v, k)
			if err != nil {
				return "", err
			}
			andFilters = append(andFilters, strings.Join(filter, " AND "))
		}
		andFilter = strings.Join(andFilters, " AND ")
		if filterObj.Or == nil {
			return andFilter, nil
		}
	}
	if filterObj.Or != nil {
		for k, v := range filterObj.Or {
			filter, err := buildFilterString(v, k)
			if err != nil {
				return "", err
			}
			orFiltersTmp := fmt.Sprintf("(%v)", strings.Join(filter, " OR "))
			orFilters = append(orFilters, orFiltersTmp)
		}
		orFilter = strings.Join(orFilters, " AND ")
		if filterObj.And == nil {
			return orFilter, nil
		}
	}
	return fmt.Sprintf("WHERE (%v) AND %v", andFilter, orFilter), nil

}

func buildFilterString(filter []FieldFilter, fieldName string) ([]string, error) {
	qfilters := []string{}
	for _, i := range filter {
		var searchParameter, operator, parseValue string
		if i.ValueType == "string" {
			parseValue = fmt.Sprintf("'%v'", i.Value)
		} else {
			parseValue = i.Value
		}
		switch i.Operator {
		case "eq":
			operator = "="
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "gt":
			operator = ">"
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "gte":
			operator = ">="
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "lt":
			operator = "<"
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "lte":
			operator = "<="
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "like":
			operator = "LIKE"
			searchParameter = fmt.Sprintf("%v %v %v", fieldName, operator, parseValue)
		case "not":
			searchParameter = fmt.Sprintf("NOT %v=%v", fieldName, parseValue)
		case "isEmpty":
			operator = "IS NULL"
			searchParameter = fmt.Sprintf("%v %v", fieldName, operator)
		case "isNotEmpty":
			operator = "IS NOT NULL"
			searchParameter = fmt.Sprintf("%v %v", fieldName, operator)
		default:
			err := errors.New("Invalid parameter")
			return nil, err
		}
		qfilters = append(qfilters, searchParameter)
	}

	return qfilters, nil
}

func prepareStatements(query string) *sql.Stmt {
	stmt, err := Connection.Prepare(query)
	if err != nil {
		utils.LogError("db", "Prepare Statement", "Fail to prepare statement")
	}
	return stmt
}

func GetRowById(table string, id int) (map[string]any, error) {
	var result = map[string]any{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=%v;", table, id)
	rows, err := Connection.Query(query)
	if err != nil {
		utils.CheckError(err, "db", "retrieve db", err.Error())
		return nil, err
	}
	col, err := rows.Columns()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	row := make([][]byte, len(col))
	rowPtr := make([]any, len(col))
	if err != nil {
		utils.CheckError(err, "db", "Retrieve column", err.Error())
		return nil, err
	}
	for i := range row {
		rowPtr[i] = &row[i]
	}
	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			utils.CheckError(err, "db", "Retrieve data", "error during data retrieval")
			return nil, err
		}
		result = getData(row, col, colTypes)
	}
	return result, nil
}

func GetRowsAll(table string) ([]map[string]any, error) {
	var result = []map[string]any{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %v;", table)
	rows, err := Connection.Query(query)
	if err != nil {
		utils.CheckError(err, "db", "retrieve db", err.Error())
		return nil, err
	}
	col, err := rows.Columns()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	row := make([][]byte, len(col))
	rowPtr := make([]any, len(col))
	if err != nil {
		utils.CheckError(err, "db", "Retrieve column", err.Error())
		return nil, err
	}
	for i := range row {
		rowPtr[i] = &row[i]
	}
	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			utils.CheckError(err, "db", "Retrieve data", "error during data retrieval")
			return nil, err
		}
		result = append(result, getData(row, col, colTypes))
	}
	return result, nil
}

func GetRowByFilter(table string, filter FilterQuery) ([]map[string]any, error) {
	var result = []map[string]any{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return nil, err
	}
	queryFilter, err := buildFilter(filter)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %v %v;", table, queryFilter)
	rows, err := Connection.Query(query)
	if err != nil {
		utils.CheckError(err, "db", "retrieve db", err.Error())
		return nil, err
	}
	col, err := rows.Columns()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		utils.CheckError(err, "db", "retrieve columns", err.Error())
		return nil, err
	}
	row := make([][]byte, len(col))
	rowPtr := make([]any, len(col))
	if err != nil {
		utils.CheckError(err, "db", "Retrieve column", err.Error())
		return nil, err
	}
	for i := range row {
		rowPtr[i] = &row[i]
	}
	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			utils.CheckError(err, "db", "Retrieve data", "error during data retrieval")
			return nil, err
		}
		result = append(result, getData(row, col, colTypes))
	}
	return result, nil
}

func InsertData(table string, inputData map[string]any) (int, error) {
	var fields, values []string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return 0, err
	}

	for k, v := range inputData {
		fields = append(fields, k)
		var parseValue string
		if typeVal := reflect.TypeOf(v).Kind(); typeVal == reflect.String {
			parseValue = fmt.Sprintf("'%v'", v)
		} else {
			parseValue = fmt.Sprintf("%v", v)
		}
		values = append(values, parseValue)
	}

	inputFields := strings.Join(fields, ",")
	inputValues := strings.Join(values, ",")
	insertStmt := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v);", table, inputFields, inputValues)
	res, err := Connection.Exec(insertStmt)
	if err != nil {
		utils.CheckError(err, "db", "insert data to db", err.Error())
		return 0, err
	}
	tmpInt, err := res.LastInsertId()
	if err != nil {
		utils.CheckError(err, "db", "Converting Id", err.Error())
		return 0, err
	}

	return int(tmpInt), nil

}

func InsertMultipleData(table string, inputDatas []map[string]any) ([]int, error) {
	var fields, vPlaceHolder []string
	var result []int

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return result, err
	}

	for k := range inputDatas[0] {
		fields = append(fields, k)
		vPlaceHolder = append(vPlaceHolder, "?")
	}
	inputFields := strings.Join(fields, ",")
	vFormats := strings.Join(vPlaceHolder, ",")

	insertStmt := fmt.Sprintf("INSERT INTO %v (%v)  VALUES (%v);", table, inputFields, vFormats)
	stmt, err := Connection.Prepare(insertStmt)

	if err != nil {
		utils.LogError("db", "Query Prep", "Error during Query preparation")
		return result, err
	}

	defer stmt.Close()

	for idx := 0; idx < len(inputDatas); idx++ {
		inputData := inputDatas[idx]
		var values []any
		for _, k := range fields {
			v := inputData[k]
			var parseValue string
			if typeVal := reflect.TypeOf(v).Kind(); typeVal == reflect.String {
				parseValue = fmt.Sprintf("'%v'", v)
			} else {
				parseValue = fmt.Sprintf("%v", v)
			}
			values = append(values, parseValue)
		}
		res, err := stmt.Exec(values...)
		if err != nil {
			utils.LogError("db", "Query Prep", "Error during insert execution")
			return result, err
		}
		var id int64
		id, err = res.LastInsertId()
		if err != nil {
			utils.LogError("db", "Query Prep", "Error during id retrieval")
			return result, err
		}
		result = append(result, int(id))
	}

	return result, nil

}

func UpdateData(table string, inputData map[string]any, id int) error {
	// This function use int Id as parameter. Change according your own requirements
	var newValues []string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return err
	}

	for k, v := range inputData {
		if k == "id" {
			continue
		}
		var parseValue string
		if typeVal := reflect.TypeOf(v).Kind(); typeVal == reflect.String {
			parseValue = fmt.Sprintf("'%v'", v)
		} else {
			parseValue = fmt.Sprintf("%v", v)
		}
		newValues = append(newValues, fmt.Sprintf("%v=%v", k, parseValue))
	}

	inputValues := strings.Join(newValues, ",")
	filter := fmt.Sprintf("id=%v", id)

	query := fmt.Sprintf("UPDATE %v SET %v WHERE %v;", table, inputValues, filter)
	_, err = Connection.Exec(query)
	if err != nil {
		utils.CheckError(err, "db", "update db", err.Error())
		return err
	}
	return nil

}

func DeleteData(table string, id int) error {
	// This function use int Id as parameter. Change according your own requirements

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return err
	}
	filter := fmt.Sprintf("id=%v", id)

	query := fmt.Sprintf("DELETE FROM %v  WHERE %v;", table, filter)
	_, err = Connection.Exec(query)
	if err != nil {
		utils.CheckError(err, "db", "delete data from db", err.Error())
		return err
	}
	return nil

}

func DeleteMultipleData(table string, id []int) error {
	// This function use int Id as parameter. Change according your own requirements

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Connection.PingContext(ctx)
	if err != nil {
		utils.CheckError(err, "db", "Database Ping", err.Error())
		return err
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %v  WHERE id=?", table)

	stmt, err := Connection.Prepare(deleteQuery)
	if err != nil {
		utils.CheckError(err, "db", "prepare delete query", err.Error())
		return err
	}

	defer stmt.Close()

	for _, i := range id {
		_, err := stmt.Exec(i)
		if err != nil {
			utils.CheckError(err, "db", "delete data from db", err.Error())
			return err
		}
	}
	return nil

}

func getData(row [][]byte, colNames []string, colTypes []*sql.ColumnType) map[string]any {
	result := map[string]any{}
	for i := 0; i < len(colNames); i++ {
		result[colNames[i]] = convert(row[i], colTypes[i])
	}
	return result
}

func convert(val []byte, colType *sql.ColumnType) any {
	var result any
	typeName := colType.DatabaseTypeName()
	switch true {
	case utils.Contains(typeName, "INT"):
		result, _ = strconv.Atoi(string(val))
	case utils.Contains(typeName, "VARCHAR"):
		result = string(val)
	case utils.Contains(typeName, "TEXT"):
		result = string(val)
	case utils.Contains(typeName, "TIMESTAMP"):
		result = string(val)
	}
	return result
}
