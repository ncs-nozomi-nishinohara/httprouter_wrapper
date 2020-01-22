package httprouter_wrapper

import (
	"net/http"
)

type Router struct{}

type ReadMe struct {
	Write     bool
	Filename  string
	Refarence bool
}

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
