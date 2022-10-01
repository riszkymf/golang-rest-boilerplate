package routes

import (
	route "github.com/riszkymf/golang-rest-boilerplate/internal/route"

	restful "github.com/emicklei/go-restful/v3"
)

func SetRoutes() {
	// Setting routes for restful endpoint, imported from route package.
	restful.Add(route.HealthRoute())
	restful.Add(route.BooksRoute())
	restful.Add(route.AuthorRoute())
	restful.Add(route.MembersRoute())
	restful.Add(route.RecordsRoute())
	restful.Add(route.RentRoute())
}
