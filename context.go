package gii

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type H map[string]interface{}

type Context struct {
	// http原始数据
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求的路径
	Path string
	// 请求的方法
	Method string
	// 路由参数
	Params map[string]string
	// http code
	StatusCode int
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{Writer: writer, Req: req, Path: req.URL.Path, Method: req.Method}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) QueryInt(key string) int {
	i, _ := strconv.Atoi(c.Req.URL.Query().Get(key))
	return i
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	// 先Header().Set
	c.SetHeader("Content-Type", "text/plain")
	// 然后WriteHeader()
	c.Status(code)
	// 最后是Write()
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
