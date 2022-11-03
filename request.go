package gorest

import (
	"encoding/json"
	"github.com/dimonrus/porterr"
	"io"
	"net/http"
)

// ParseJsonBody Read json data from request body
func ParseJsonBody(r io.ReadCloser, data interface{}) porterr.IError {
	if r == http.NoBody {
		return nil
	}
	body, err := io.ReadAll(r)
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
