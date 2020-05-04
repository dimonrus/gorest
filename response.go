package gorest

import (
	"encoding/json"
	"github.com/dimonrus/porterr"
	"io/ioutil"
	"net/http"
)

// New Json Response without error
func NewOkJsonResponse(message interface{}, data interface{}, meta interface{}) *JsonResponse {
	return &JsonResponse{HttpCode: http.StatusOK, Message: message, Data: data, Meta: meta}
}

// New Json Response with error
func NewErrorJsonResponse(e porterr.IError) *JsonResponse {
	httpCode := http.StatusInternalServerError
	if e != nil && e.GetHTTP() >= http.StatusBadRequest && e.GetHTTP() <= http.StatusNetworkAuthenticationRequired {
		httpCode = e.GetHTTP()
	}
	return &JsonResponse{HttpCode: httpCode, Error: e}
}

// Send response to client
func Send(writer http.ResponseWriter, response *JsonResponse) {
	SendJson(writer, response.HttpCode, response)
}

// Sent json into http writer
func SendJson(writer http.ResponseWriter, httpCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpCode)
	if data == nil {
		return
	}
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

// Default response strategy
func ResponseErrorStrategy(response *http.Response) porterr.IError {
	var e porterr.IError
	if response.StatusCode >= http.StatusBadRequest {
		// Default error message
		message := http.StatusText(response.StatusCode) + ": " + response.Request.URL.Path + " Service: " + response.Request.Host
		// Init Error
		e = porterr.New(porterr.PortErrorResponse, message).HTTP(response.StatusCode)
		// Init response struct
		result := &JsonResponse{Error: e}
		// Read response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return e.PushDetail(porterr.PortErrorBody, "body", err.Error())
		}
		if len(body) != 0 {
			err = json.Unmarshal(body, result)
			if err != nil {
				e = e.PushDetail(porterr.PortErrorBody, "body", err.Error())
			}
		}
	}
	return e
}
