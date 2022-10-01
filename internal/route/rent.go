package route

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type RentData struct {
	Id         int    `json:"id" default:"-1"`
	BookId     int    `json:"book_id" default:"-1"`
	MemberId   int    `json:"member_id" default:"-1"`
	AuthorName string `json:"author_name" default:""`
	Email      string `json:"email" default:""`
	Firstname  string `json:"firstname" default:""`
	Lastname   string `json:"lastname" default:""`
	RentDate   string `json:"rent_date" default:""`
	DueDate    string `json:"due_date" default:""`
	RentStatus string `json:"rent_status" default:""`
}

func RentRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/rent").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/{records-id}").
		To(GetRentData)).
		Doc("Retrieve rent by ID").
		Param(service.PathParameter("record-id", "Identifier of record").DataType("integer"))
	service.Route(service.GET("/").
		To(GetAllRentData)).
		Doc("Retrieve available rent data")
	return service
}

func GetAllRentData(request *restful.Request, response *restful.Response) {
	data, err := handler.GetRowsAll("v_rent")
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: data, Errors: nil, StatusCode: 200}
	response.WriteAsJson(res)
}

func GetRentData(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("record-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("GetRecord", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	data, err := handler.GetRowById("records", idParse)
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: data, StatusCode: http.StatusOK}
	response.WriteAsJson(res)
}
