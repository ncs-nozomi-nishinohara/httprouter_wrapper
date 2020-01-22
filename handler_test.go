package httprouter_wrapper_test

import (
	"encoding/json"
	"fmt"
	"github.com/ncs-nozomi-nishinohara/httprouter_wrapper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
}

// レシーバーの定義
func (Router) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var js = `{"result": "true"}`
	var result_interface interface{}
	json.Unmarshal([]byte(js), &result_interface)
	result, _ := json.Marshal(result_interface)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// レシーバーの定義
func (Router) Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var js = `{"result": "true"}`
	var result_interface interface{}
	json.Unmarshal([]byte(js), &result_interface)
	result, _ := json.Marshal(result_interface)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (Router) AllMethod(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var js = `{"result": "true"}`
	var result_interface interface{}
	json.Unmarshal([]byte(js), &result_interface)
	result, _ := json.Marshal(result_interface)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	switch r.Method {
	case http.MethodGet:
		w.Write(result)
	case http.MethodPost:
		w.Write(result)
	case http.MethodPut:
		w.Write(result)
	case http.MethodDelete:
		w.Write(result)
	}
}

func Request(r *http.Request, err error) *http.Request {
	return r
}
func TestNew(t *testing.T) {
	var test_json = `{"result": "true"}`
	var test_interface interface{}
	wr := Router{}
	json.Unmarshal([]byte(test_json), &test_interface)
	result, _ := json.Marshal(test_interface)
	router := httprouter_wrapper.NewRouterWrapperHandler("test_service.yaml", httprouter_wrapper.ReadMe{
		Write:    false,
		Filename: "TestReadme.md",
	})
	router.Router = wr
	httprouter_wrapper.New(router)
	var requests []*http.Request

	requests = append(requests, Request(http.NewRequest(http.MethodGet, "/test", nil)))
	requests = append(requests, Request(http.NewRequest(http.MethodPost, "/test", nil)))

	for _, method := range []string{http.MethodDelete, http.MethodGet, http.MethodPost} {
		requests = append(requests, Request(http.NewRequest(method, "/all_method", nil)))
	}

	for _, req := range requests {
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)
		var response_interface interface{}
		json.Unmarshal([]byte(w.Body.Bytes()), &response_interface)
		response, _ := json.Marshal(response_interface)
		if w.Result().StatusCode == 200 {
			if string(result) != string(response) {
				t.Fatalf("Response Data Error. %v", w.Body.String())
			}
		} else {
			t.Fatalf("Status Code Error %v", w.Result().StatusCode)
		}
	}
	r := Request(http.NewRequest(http.MethodGet, "/refarence", nil))
	w := httptest.NewRecorder()
	router.Handler.ServeHTTP(w, r)
	if w.Result().StatusCode == 200 {
		fmt.Println(w.Body.String())
	}

}
