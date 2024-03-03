package response

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

type ApiError struct {
	status   int
	ErrorMsg string `json:"error"`
}

func NewApiError(status int, err error) *ApiError {
	_, file, line, _ := runtime.Caller(1)
	log.Println(file, line, err)
	return &ApiError{
		status,
		err.Error(),
	}
}

func (e *ApiError) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(e.status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Fatal(err.Error())
	}
}

//// TODO: maybe change render to something else, since we don't really need to handle errors
//func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
//	render.Status(r, e.status)
//}
