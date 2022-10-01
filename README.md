# golang-rest-boilerplate
Golang Rest API Boilerplate

# Description
golang-rest-boilerplate is a REST API Boilerplate written on golang using go-restful package and using sqlite3 as database. This Boilerplate is also a simulation for a library system as an example for other endpoint routes.

# File structure
```
./
├─ cmd/
│  ├─ main.go
├─ internal/
│  ├─ src/
│  │  ├─ utils.go
│  ├─ handler/
│  │  ├─ dbHandler.go
│  ├─ route/
│  │  ├─ yourRoute.go
│  ├─ routes.go
├─ .gitignore
├─ package.json
├─ README.md
```

A new endpoint handler must be written under route directory where routes.go will import that endpoint route, which will be used on main.go

### Example

Create your new route on : 
internal/route/aNewRoute.go
```go
package route

import (
	"fmt"
	restful "github.com/emicklei/go-restful/v3"
)

type NewRoute struct {
	Message    string `json:"message"`
}

func NewRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/newroute").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/").
		To(GetMyRoute)).
		Doc("Check out my new route")
	return service
}

func GetMyRoute(request *restful.Request, response *restful.Response) {
	res := NewRoute{Message: "OK"0}
	response.WriteAsJson(res)
	return
}
```
And modify `internal/routes.go` content with your NewRoute function
```go
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
	restful.Add(route.NewRoute()) // your new route
}

```
## dbHandler Filtering
Filtering on this boilerplate is build on this structure
```go
type FilterQuery struct {
	And map[string][]FieldFilter

	Or map[string][]FieldFilter
}

type FieldFilter struct {
	Operator  string
	Value     string
	ValueType string
}
```
Where possible operator values are:
    - eq: equal; columnName = value
    - gt: greater than; columnName > value
    - gte: greater than equal; columnName >= value
    - lt: less than; columnName < value
    - lte: less than equal; columnName<= value
    - like: Like; columnName LIKE value
    - not: Not; NOT columnName=value
    - isEmpty: is null; columnName IS NULL
    - isNotEmpty: not null; columnName IS NOT NULL

All values on FilterQuery.And must be fulfilled where only one value on each fields in FilterQuery.Or is have to be fulfilled.

For example:

| id | title                     | author_name       | stock |
|----|---------------------------|-------------------|-------|
| 1  | Discworld: Guards!Guards! | Terry Pratchett   | 10    |
| 2  | Discworld: Nightwatch     | Terry Pratchett   | 5     |
| 3  | Dubliners                 | James Joyce       | 7     |
| 4  | The Carpet People         | Terry Pratchett   | 6     |
| 5  | Finnegans' Wake           | James Joyce       | 5     |
| 6  | Moby Dick                 | Herman Melville   | 15    |
| 7  | Citizen Kane              | Herman Mankiewicz | 9     |

Case : We need to query for any book from Terry Pratchett with stock below 10 (10 not included).
Then the filter will be
```go
	myFilter := FilterQuery{
	And: map[string][]FieldFilter{
		"author_name": {
			{
				Operator:  "like",
				Value:     "Terry Pratchett",
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
}
```

This will generate the following filter on query
```sql
WHERE (author_name LIKE 'Terry Pratchett' AND stock < 10)
```

What if we want to query for books by James Joyce or any author with Herman as firstname and stock above 6 ?
```go
	myFilter := FilterQuery{
		And: map[string][]handler.FieldFilter{
			"stock": {
				{
					Operator:  "gt",
					Value:     "6",
					ValueType: "int",
				},
			},
		},
		Or: map[string][]handler.FieldFilter{
			"author_name": {
				{
					Operator:  "like",
					Value:     "Herman%",
					ValueType: "string",
				},
				{
					Operator:  "like",
					Value:     "James Joyce",
					ValueType: "string",
				},
			},
		},
	}
```
Filter above will generate the following sql query
```WHERE (stock > 6) AND (author_name LIKE 'Herman%' OR author_name LIKE 'James Joyce')```