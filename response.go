package gorest

import (
	"encoding/json"
	"github.com/dimonrus/porterr"
	"net/http"
)

// New Json Response without error
func NewOkJsonResponse(message interface{}, data interface{}, meta interface{}) *JsonResponse {
	return &JsonResponse{HttpCode: http.StatusOK, Message: message, Data: data, Meta: meta}
}

// New Json Response with error
func NewErrorJsonResponse(e porterr.IError, codeMap HttpCodeMap) *JsonResponse {
	httpCode := http.StatusInternalServerError
	if err, ok := e.(*porterr.PortError); ok {
		if v, ok := err.Code.(string); ok {
			if code, ok := codeMap[v]; ok {
				httpCode = code
			}
		} else if v, ok := err.Code.(int); ok {
			if v >= http.StatusBadRequest && v <= http.StatusNetworkAuthenticationRequired {
				httpCode = v
			}
		}
	}
	return &JsonResponse{Message: e.Error(), HttpCode: httpCode, Error: e}
}

// Send response to client
func Send(writer http.ResponseWriter, response *JsonResponse) {
	SendJson(writer, response.HttpCode, response)
}

// Sent json into http writer
func SendJson(writer http.ResponseWriter, httpCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpCode)
	body, err := json.Marshal(data)
	if err != nil {
		_, err := writer.Write([]byte("JSON marshal failed: " + err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}
	_, err = writer.Write(body)
	if err != nil {
		panic(err)
	}
}
