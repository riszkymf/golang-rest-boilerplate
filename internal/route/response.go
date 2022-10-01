package route

type ResponseObj struct {
	Data       interface{} `json:"data"`
	Errors     []string    `json:"error"`
	StatusCode int         `json:"status"`
}
