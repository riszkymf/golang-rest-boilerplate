package route

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type MemberResponse struct {
	Data       interface{} `json:"data"`
	Errors     []string    `json:"error"`
	StatusCode int         `json:"status"`
}

type Members struct {
	Id        int    `json:"id" default:"-1"`
	Firstname string `json:"firstname" default:""`
	Lastname  string `json:"lastname" default:""`
	Email     string `json:"email" default:""`
	Address   string `json:"address" default:""`
}

func MembersRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/members").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/{member-id}").
		To(GetMember)).
		Doc("Retrieve member by ID").
		Param(service.PathParameter("member-id", "Identifier of member").DataType("integer"))
	service.Route(service.GET("/").
		To(GetAllMembers)).
		Doc("Retrieve available members")
	service.Route(service.POST("").
		To(InsertMember)).
		Doc("Insert new member")
	service.Route(service.POST("/{member-id}").
		To(UpdateMember)).
		Doc("Update member by ID").
		Param(service.PathParameter("member-id", "Identifier of member").DataType("integer"))
	return service
}

func GetAllMembers(request *restful.Request, response *restful.Response) {
	member, err := handler.GetRowsAll("members")
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: member, Errors: nil, StatusCode: 200}
	response.WriteAsJson(res)
}

func GetMember(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("member-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("GetMember", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	member, err := handler.GetRowById("members", idParse)
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: member, StatusCode: http.StatusOK}
	response.WriteAsJson(res)
}

func InsertMember(request *restful.Request, response *restful.Response) {
	member := new(Members)
	err := request.ReadEntity(&member)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest}
		response.WriteAsJson(res)
		return
	}

	inputData := map[string]any{
		"firstname": member.Firstname,
		"lastname":  member.Lastname,
		"email":     member.Email,
		"address":   member.Address,
	}
	addMember, err := handler.InsertData("members", inputData)
	if err != nil {
		res := ResponseObj{
			Errors:     []string{err.Error()},
			StatusCode: http.StatusInternalServerError,
		}
		response.WriteAsJson(res)
	}
	member.Id = addMember
	response.WriteAsJson(ResponseObj{Data: member})

}

func UpdateMember(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("member-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Update Member", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}

	var updateInput map[string]interface{}

	member := Members{
		Id:        idParse,
		Address:   "",
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}

	err = request.ReadEntity(&updateInput)

	filteredInput, err := utils.FilterInputMap(member, updateInput)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}

	err = handler.UpdateData("member", filteredInput, member.Id)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	filteredInput["id"] = member.Id
	response.WriteAsJson(ResponseObj{Data: filteredInput})

}

func DeleteMember(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("member-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("DeleteMember", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	err = handler.DeleteData("members", idParse)
	if err != nil {
		response.WriteAsJson(ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest})
		return
	}
}
