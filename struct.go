package httprouter_wrapper

import (
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
