package routes

import (
	"fmt"

	route "github.com/riszkymf/golang-rest-boilerplate/internal/route"

	restful "github.com/emicklei/go-restful/v3"
	utils "github.com/riszkymf/golang-rest-boilerplate/internal/src"
)

type RouteFilterConfig struct {
	WebServiceLogging string
	Auth              string
}

func SetRoutes(routeContainer *restful.Container) *restful.Container {
	// Setting routes for restful endpoint, imported from route package.

	routeContainer.Add(route.HealthRoute())
	routeContainer.Add(route.BooksRoute())
	routeContainer.Add(route.AuthorRoute())
	routeContainer.Add(route.MembersRoute())
	routeContainer.Add(route.RecordsRoute())
	routeContainer.Add(route.RentRoute())
	return routeContainer

}

func SetFilters(routeContainer *restful.Container, config RouteFilterConfig) *restful.Container {
	/*
		Config Model:
		type RouteFilterConfig struct{
			webServiceLogging string
			Authorization	  string
			Authentication    string
		}
	*/

	if config.WebServiceLogging == "TRUE" {
		utils.LogInfo("[webservice-init]", "initalizing filter", "adding logging to filters")
		routeContainer.Filter(webserviceLogging)
	}

	return routeContainer
}

func webserviceLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	logMsg := fmt.Sprintf("%s,%s,%s\n", req.Request.Method, req.Request.URL, req.Request.RemoteAddr)
	utils.LogInfo("[webservice-logging] ", "log", logMsg)
	chain.ProcessFilter(req, resp)
}
