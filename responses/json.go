package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errorString string) {
	if statusCode >= 499 {
		log.Printf("Internal server error : %v", errorString)
		return
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	errorRes := errorResponse{
		Error: errorString,
	}
	RespondWithJSON(w, statusCode, errorRes)
}
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload int json : %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}