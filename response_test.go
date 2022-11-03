package gorest

import (
	"encoding/json"
	"fmt"
	"github.com/dimonrus/porterr"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testDataStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewOkJsonResponse(t *testing.T) {
	data := testDataStruct{Id: 10, Name: "test string"}
	meta := Meta{Page: 10, Limit: 1, Total: 10}
	resp := NewOkJsonResponse("Success", data, &meta)
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%s", body) != `{"message":"Success","data":{"id":10,"name":"test string"},"meta":{"page":10,"limit":1,"total":10}}` {
		t.Fatal("wrong format")
	}
}

func TestNewErrorJsonResponse(t *testing.T) {
	e := porterr.New(porterr.PortErrorSystem, "Some failed")
	e = e.PushDetail("SOME_ERROR", "failed", "First error")
	e = e.PushDetail("SOME_ERROR", "", "Second error")

	body, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s", body)
	if fmt.Sprintf("%s", body) != `{"code":"PORTABLE_ERROR_SYSTEM","message":"Some failed","data":[{"code":"SOME_ERROR","name":"failed","message":"First error"},{"code":"SOME_ERROR","message":"Second error"}]}` {
		t.Fatal("wrong format")
	}
}

func testErrorHandler(w http.ResponseWriter, r *http.Request) {
	Send(w, NewErrorJsonResponse(porterr.New(porterr.PortErrorRequest, "Some failed message").HTTP(http.StatusBadRequest)))
}
func testErrorMapHandler(w http.ResponseWriter, r *http.Request) {
	e := porterr.New(porterr.PortErrorSearch, "Some failed message").HTTP(http.StatusNotFound)
	e = e.PushDetail(porterr.PortErrorDecoder, "", "Some error")

	Send(w, NewErrorJsonResponse(e))
}

func TestBadResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	testErrorHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("wrong status")
	}
	res := struct {
		Error porterr.IError
	}{
		Error: &porterr.PortError{},
	}
	err := ParseJsonBody(resp.Body, &res)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(res.Error)
}

func TestNotFoundError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	testErrorMapHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("wrong status")
	}
	res := struct {
		Error porterr.IError
	}{
		Error: &porterr.PortError{},
	}
	err := ParseJsonBody(resp.Body, &res)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Error.GetDetails()) != 1 {
		t.Fatal("wrong error detail length")
	}
}
