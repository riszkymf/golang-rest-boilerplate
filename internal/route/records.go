package route

import (
	"net/http"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type Records struct {
	Id         int    `json:"id" default:"-1"`
	BookId     int    `json:"book_id" default:"-1"`
	MemberId   int    `json:"member_id" default:"-1"`
	RentDate   string `json:"rent_date" default:""`
	DueDate    string `json:"due_date" default:""`
	RentStatus string `json:"rent_status" default:""`
}

func RecordsRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/records").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/{records-id}").
		To(GetRecord)).
		Doc("Retrieve record by ID").
		Param(service.PathParameter("record-id", "Identifier of record").DataType("integer"))
	service.Route(service.GET("/").
		To(GetAllRecords)).
		Doc("Retrieve available records")
	service.Route(service.POST("").
		To(InsertRecord)).
		Doc("Insert new record")
	service.Route(service.POST("/{record-id}").
		To(UpdateRecord)).
		Doc("Update record by ID").
		Param(service.PathParameter("record-id", "Identifier of record").DataType("integer"))
	return service
}

func GetAllRecords(request *restful.Request, response *restful.Response) {
	records, err := handler.GetRowsAll("records")
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: records, Errors: nil, StatusCode: 200}
	response.WriteAsJson(res)
}

func GetRecord(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("record-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("GetRecord", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	record, err := handler.GetRowById("records", idParse)
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: record, StatusCode: http.StatusOK}
	response.WriteAsJson(res)
}

func InsertRecord(request *restful.Request, response *restful.Response) {
	record := new(Records)
	err := request.ReadEntity(&record)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest}
		response.WriteAsJson(res)
		return
	}

	inputData := map[string]any{
		"member_id":   record.MemberId,
		"book_id":     record.BookId,
		"rent_date":   record.RentDate,
		"due_date":    record.DueDate,
		"rent_status": record.RentStatus,
	}
	addRecord, err := handler.InsertData("records", inputData)
	if err != nil {
		res := ResponseObj{
			Errors:     []string{err.Error()},
			StatusCode: http.StatusInternalServerError,
		}
		response.WriteAsJson(res)
	}
	record.Id = addRecord
	response.WriteAsJson(ResponseObj{Data: record})

}

func UpdateRecord(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("record-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Update Record", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}

	var updateInput map[string]interface{}

	record := Records{
		Id:         idParse,
		MemberId:   -1,
		BookId:     -1,
		DueDate:    time.Unix(0, 0).Format("2022-06-19"),
		RentDate:   time.Unix(0, 0).Format("2022-06-19"),
		RentStatus: "",
	}

	err = request.ReadEntity(&updateInput)

	filteredInput, err := utils.FilterInputMap(record, updateInput)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}

	err = handler.UpdateData("records", filteredInput, record.Id)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	filteredInput["id"] = record.Id
	response.WriteAsJson(ResponseObj{Data: filteredInput})

}

func DeleteRecord(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("record-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Delete Records", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	err = handler.DeleteData("records", idParse)
	if err != nil {
		response.WriteAsJson(ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest})
		return
	}
}
