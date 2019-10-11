package gorest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestJsonRequest struct {
	Number int `json:"number"`
}

func TestReadJsonBody(t *testing.T) {
	n := TestJsonRequest{123}
	r := &TestJsonRequest{}
	b, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/no_url", bytes.NewReader(b))
	e := ParseJsonBody(req.Body, r)
	if e != nil {
		t.Fatal(e.Error())
	}
	if r.Number != 123 {
		t.Fatal("incorrect number")
	}
}

type testResponse struct {
	Message string
	Data    interface{}
}

type okStruct struct {
	Id     int
	Custom string
}

func testOkHandler(w http.ResponseWriter, r *http.Request) {
	os := okStruct{
		Id:     12,
		Custom: "message",
	}

	Send(w, NewOkJsonResponse("Cool its works", os, nil))
}

func TestOkResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	testOkHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("wrong status")
	}
	os := okStruct{}
	rsp := testResponse{Data: &os}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		t.Fatal(err)
	}
	if os.Id != 12 {
		t.Fatal("wrong id")
	}
}