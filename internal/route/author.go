package route

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type AuthorResponse struct {
	Data       interface{} `json:"data"`
	Errors     []string    `json:"error"`
	StatusCode int         `json:"status"`
}

type Author struct {
	Id   int    `json:"id" default:"-1"`
	Name string `json:"title" default:""`
}

func AuthorRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/author").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/{author-id}").
		To(GetAuthor)).
		Doc("Retrieve author by ID").
		Param(service.PathParameter("author-id", "Identifier of author").DataType("integer"))
	service.Route(service.GET("/").
		To(GetAllAuthors)).
		Doc("Retrieve available authors")
	service.Route(service.POST("").
		To(InsertAuthor)).
		Doc("Insert new author")
	service.Route(service.POST("/{author-id}").
		To(UpdateAuthor)).
		Doc("Update author by ID").
		Param(service.PathParameter("author-id", "Identifier of author").DataType("integer"))
	return service
}

func GetAllAuthors(request *restful.Request, response *restful.Response) {
	author, err := handler.GetRowsAll("author")
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: author, Errors: nil, StatusCode: 200}
	response.WriteAsJson(res)
}

func GetAuthor(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("author-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("GetAuthor", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	author, err := handler.GetRowById("author", idParse)
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: author, StatusCode: http.StatusOK}
	response.WriteAsJson(res)
}

func InsertAuthor(request *restful.Request, response *restful.Response) {
	author := new(Author)
	err := request.ReadEntity(&author)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest}
		response.WriteAsJson(res)
		return
	}

	inputData := map[string]any{
		"name": author.Name,
	}
	addAuthor, err := handler.InsertData("author", inputData)
	if err != nil {
		res := ResponseObj{
			Errors:     []string{err.Error()},
			StatusCode: http.StatusInternalServerError,
		}
		response.WriteAsJson(res)
	}
	author.Id = addAuthor
	response.WriteAsJson(ResponseObj{Data: author})

}

func UpdateAuthor(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("author-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Update Author", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}

	var updateInput map[string]interface{}

	author := Author{
		Id:   idParse,
		Name: "",
	}

	err = request.ReadEntity(&updateInput)

	filteredInput, err := utils.FilterInputMap(author, updateInput)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}

	err = handler.UpdateData("author", filteredInput, author.Id)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	filteredInput["id"] = author.Id
	response.WriteAsJson(ResponseObj{Data: filteredInput})

}

func DeleteAuthor(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("author-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("DeleteAuthor", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	err = handler.DeleteData("author", idParse)
	if err != nil {
		response.WriteAsJson(ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest})
		return
	}
}
