package main

import (
	"gii"
	"log"
	"net/http"
)

func main() {
	engine := gii.New()

	engine.GET("/", func(c *gii.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gii</h1>")
	})

	// get
	// http://localhost:8080/hello?name=gii&age=18
	engine.GET("/hello", func(c *gii.Context) {
		c.String(http.StatusOK, "hello %s, %d", c.Query("name"), c.QueryInt("age"))
	})

	// get
	// http://localhost:8080/hello/gii
	engine.GET("/hello/:name", func(ctx *gii.Context) {
		ctx.String(http.StatusOK, "hello %s", ctx.Param("name"))
	})

	engine.GET("/assert/*filepath", func(ctx *gii.Context) {
		ctx.JSON(http.StatusOK, gii.H{
			"filepath": ctx.Param("filepath"),
		})
	})

	// post
	// http://localhost:8080/login
	// x-www-form-urlencoded:
	// user=gii&password=123456
	engine.POST("/login", func(c *gii.Context) {
		c.JSON(http.StatusOK, gii.H{
			"username": c.PostForm("user"),
			"password": c.PostForm("password"),
		})
	})

	log.Fatal(engine.Run(":8080"))
}
