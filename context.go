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
	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{Writer: writer, Req: req, Path: req.URL.Path, Method: req.Method, index: -1}
}

// Next 执行下一个中间件
// func A(c *Context) {
//    part1
//    c.Next()
//    part2
// }
// func B(c *Context) {
//    part3
//    c.Next()
//    part4
// }
//
// 假设我们应用了中间件 A 和 B，和路由映射的 Handler。c.handlers是这样的[A, B, Handler]，c.index初始化为-1。
// 调用c.Next()，接下来的流程是这样的：
// c.index++，c.index 变为 0
// 0 < 3，调用 c.handlers[0]，即 A
// 执行 part1，调用 c.Next()
// c.index++，c.index 变为 1
// 1 < 3，调用 c.handlers[1]，即 B
// 执行 part3，调用 c.Next()
// c.index++，c.index 变为 2
// 2 < 3，调用 c.handlers[2]，即Handler
// Handler 调用完毕，此时c.index = 3 循环结束， 返回到 B 中的 part4，执行 part4
// part4 执行完毕，返回到 A 中的 part2，执行 part2
// part2 执行完毕，结束。
// 最终的顺序是part1 -> part3 -> Handler -> part 4 -> part2
func (c *Context) Next() {
	// index是记录当前执行到第几个中间件，当在中间件中调用Next方法时，控制权交给了下一个中间件，
	// 直到调用到最后一个中间件，然后再从后往前，调用每个中间件在Next方法之后定义的部分
	c.index++
	s := len(c.handlers) // 中间件和用户自定义的handler外是同一个类型
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c Context) Abort(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"code": code, "message": err})
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
