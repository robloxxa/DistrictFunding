package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

// Error function responses with an error message in json format
// Panics if encoding of json is failed
func Error(w http.ResponseWriter, status int, errMsg error) {
	w.WriteHeader(status)
	Json(w, &response{errMsg.Error()})
}

func Message(w http.ResponseWriter, msg string) {
	Json(w, &response{msg})
}

func Json(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

//// TODO: maybe change render to something else, since we don't really need to handle errors
//func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
//	render.Status(r, e.status)
//}
