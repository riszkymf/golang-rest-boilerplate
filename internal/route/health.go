package route

import (
	restful "github.com/emicklei/go-restful/v3"
)

type Health struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

func HealthRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/health").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/").
		To(GetUser)).
		Doc("Health Check")
	return service
}

func GetUser(request *restful.Request, response *restful.Response) {
	res := Health{Message: "OK", StatusCode: 200}
	response.WriteAsJson(res)
	return
}
