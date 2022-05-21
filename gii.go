package gii

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gii
type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if h, ok := e.router[key]; ok {
		h(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 请求方法和请求路径作为key
	e.router[method+"-"+pattern] = handler
}

// GET defines the method to add GET request
func (e *Engine) GET(patter string, handler HandlerFunc) {
	e.addRoute("GET", patter, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(patter string, handler HandlerFunc) {
	e.addRoute("POST", patter, handler)
}

// Run defines the method to start a http server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
