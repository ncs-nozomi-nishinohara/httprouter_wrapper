package wrapper_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code        int `json:"code"`
	Process_cnt int `json:"process_cnt"`
	Error       struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Error(w http.ResponseWriter, err error, status ...interface{}) bool {
	if err != nil {
		log.Println(err)
		interfaceflg := false
		result := ErrorResponse{}
		if len(status) >= 1 {
			switch status[0].(type) {
			case int:
				result.Code = status[0].(int)
			case ErrorResponse:
				result = status[0].(ErrorResponse)
				interfaceflg = true
			}
		}
		if interfaceflg {
			jsondata, _ := json.Marshal(result)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsondata)
		} else {
			result.Process_cnt = 0
			if len(status) == 2 {
				result.Error.Code = status[1].(int)
			} else {
				result.Error.Code = 400
			}
			result.Error.Message = err.Error()
			jsondata, _ := json.Marshal(result)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsondata)
		}
		return true
	}
	return false
}
