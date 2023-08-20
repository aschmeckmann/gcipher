package api

import (
	"encoding/json"
	"net/http"
)

func EncodeResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	encodeJSONResponse(w, response, http.StatusOK)
}

func EncodeErrorResponse(w http.ResponseWriter, errorCode int, errorMessage string) {
	errors := []Error{
		{
			Code:    errorCode,
			Message: errorMessage,
		},
	}
	response := Response{
		Success: false,
		Errors:  errors,
	}
	encodeJSONResponse(w, response, errorCode)
}

func encodeJSONResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
