package gii

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gii
type HandlerFunc func(ctx *Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: 后续优化：使用对象池 sync.Pool复用对象，减少内存分配、释放和GC
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 请求方法和请求路径作为key
	e.router.addRoute(method, pattern, handler)
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
