package gorest

import (
	"encoding/json"
	"github.com/dimonrus/porterr"
	"io"
	"io/ioutil"
	"net/http"
)

// Прочитать тело дескриптор запроса
func ParseJsonBody(r io.ReadCloser, data interface{}) porterr.IError {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return porterr.New(porterr.PortErrorIO, "I/O error. Request body error: "+err.Error()).HTTP(http.StatusBadRequest)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			panic(err)
		}
	}()
	err = json.Unmarshal(body, data)
	if err != nil {
		return porterr.New(porterr.PortErrorDecoder, "Unmarshal error: "+err.Error()).HTTP(http.StatusBadRequest)
	}
	return nil
}
