package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	ErrorInternal      = "internal"
	ErrorMalformedJSON = "malformed-json"
)

type ErrorResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
func (a *Server) WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		http.Error(w, `{"error":"`+ErrorInternal+`"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		// Very unlikely to happen, but log any error (not much more we can do)
		log.Printf("error writing JSON: %v", err)
	}
}

func (a *Server) JsonError(w http.ResponseWriter, status int, error string, data map[string]interface{}) {
	response := ErrorResponse{
		Status:  status,
		Message: error,
		Data:    data,
	}
	a.WriteJSON(w, status, response)
}

func (a *Server) ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading JSON body: %v", err)
		a.JsonError(w, http.StatusInternalServerError, ErrorInternal, nil)
		return false
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		data := map[string]interface{}{"message": err.Error()}
		a.JsonError(w, http.StatusBadRequest, ErrorMalformedJSON, data)
		return false
	}
	return true
}
