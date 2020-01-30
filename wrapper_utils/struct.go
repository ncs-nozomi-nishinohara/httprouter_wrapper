package wrapper_utils

import (
	"fmt"
	"net/http"
)

// レシーバー構造体
type Router struct{}

// readme構造体
type ReadMe struct {
	Write     bool
	Filename  string
	Refarence bool
}

// wrapperhandler構造体
type RouterWrapperHandler struct {
	Filename    string
	Readme      ReadMe
	key         string
	Handler     http.Handler
	port        string
	error_      error
	ListenServe func() error
	Router      interface{}
}

// wrapperハンドラーエラーメソッド
func (w *RouterWrapperHandler) Error() string {
	if w.error_ != nil {
		return w.error_.Error()
	}
	return fmt.Sprintf("%s error key %s", w.Filename, w.key)
}

func (w *RouterWrapperHandler) SetError(err error) {
	w.error_ = err
}

func (w *RouterWrapperHandler) SetKey(key string) {
	w.key = key
}

func (w *RouterWrapperHandler) SetPort(port string) {
	w.port = ":" + port
}

func (w *RouterWrapperHandler) GetPort() string {
	return w.port
}
