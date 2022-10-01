package route

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/riszkymf/golang-rest-boilerplate/internal/handler"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type Book struct {
	Id       int    `json:"id" default:"-1"`
	Title    string `json:"title" default:""`
	AuthorId int    `json:"author_id" default:""`
	Stock    int    `json:"stock" default:"-1"`
}

func BooksRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/books").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/{book-id}").
		To(GetBook)).
		Doc("Retrieve book by ID").
		Param(service.PathParameter("book-id", "Identifier of book").DataType("integer"))
	service.Route(service.GET("/").
		To(GetAllBooks)).
		Doc("Retrieve available books")
	service.Route(service.POST("").
		To(InsertBook)).
		Doc("Insert new book")
	service.Route(service.POST("/{book-id}").
		To(UpdateBook)).
		Doc("Update book by ID").
		Param(service.PathParameter("book-id", "Identifier of book").DataType("integer")).
		Param(service.BodyParameter("body", "testing").DataType("Book"))
	return service
}

func GetAllBooks(request *restful.Request, response *restful.Response) {
	books, err := handler.GetRowsAll("v_books")
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: books, Errors: nil, StatusCode: 200}
	response.WriteAsJson(res)
}

func GetBook(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("book-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("GetBook", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	book, err := handler.GetRowById("books", idParse)
	if err != nil {
		res := ResponseObj{Data: nil, Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	res := ResponseObj{Data: book, StatusCode: http.StatusOK}
	response.WriteAsJson(res)
}

func InsertBook(request *restful.Request, response *restful.Response) {
	book := new(Book)
	err := request.ReadEntity(&book)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest}
		response.WriteAsJson(res)
		return
	}
	authorId := book.AuthorId
	author, err := handler.GetRowById("author", authorId)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	if author["id"] == nil {
		res := ResponseObj{
			Errors:     []string{"author_id does not exist"},
			StatusCode: http.StatusNoContent,
		}
		response.WriteAsJson(res)
		return
	}
	inputData := map[string]any{
		"title":     book.Title,
		"stock":     book.Stock,
		"author_id": book.AuthorId,
	}
	addBook, err := handler.InsertData("books", inputData)
	if err != nil {
		res := ResponseObj{
			Errors:     []string{err.Error()},
			StatusCode: http.StatusInternalServerError,
		}
		response.WriteAsJson(res)
	}
	book.Id = addBook
	response.WriteAsJson(ResponseObj{Data: book})

}

func UpdateBook(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("book-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Update Book", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}

	var updateInput map[string]interface{}

	book := Book{
		Id:       idParse,
		Title:    "",
		Stock:    -1,
		AuthorId: -1,
	}

	err = request.ReadEntity(&updateInput)

	filteredInput, err := utils.FilterInputMap(book, updateInput)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}

	err = handler.UpdateData("books", filteredInput, book.Id)
	if err != nil {
		res := ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusInternalServerError}
		response.WriteAsJson(res)
		return
	}
	filteredInput["id"] = book.Id
	response.WriteAsJson(ResponseObj{Data: filteredInput})

}

func DeleteBook(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("book-id")
	idParse, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError("Delete Book", "Invalid Id Type", err.Error())
		response.WriteAsJson(ResponseObj{Data: nil, Errors: []string{"ID must be numerical"}, StatusCode: http.StatusBadRequest})
		return
	}
	err = handler.DeleteData("books", idParse)
	if err != nil {
		response.WriteAsJson(ResponseObj{Errors: []string{err.Error()}, StatusCode: http.StatusBadRequest})
		return
	}
}
