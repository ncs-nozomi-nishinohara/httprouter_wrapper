package wrapper_utils

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
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
	Filename  string
	Readme    ReadMe
	Handler   http.Handler
	Router    interface{}
	Migration bool
	address   string
	port      string
	error_    error
	key       string
}

func (w *RouterWrapperHandler) ListenServe() error {
	return http.ListenAndServe(w.GetPort(), w.Handler)
}

func (w *RouterWrapperHandler) ListenFastCGI() error {
	l, _ := net.Listen("tcp", w.GetPort())
	return fcgi.Serve(l, w.Handler)
}

func (w *RouterWrapperHandler) GetAdress() string {
	return w.address
}

func (w *RouterWrapperHandler) SetAddress(address string) {
	w.address = address
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
