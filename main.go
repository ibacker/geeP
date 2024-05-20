package main

import (
	"fmt"
	"geeP/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.GET("/hello/:param", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("param"), c.Path)
	})

	r.GET("/hello2/*", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	// 添加 group 后
	r.GET("/index", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Index</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
		v1.POST("/post", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"name": c.PostForm("name")})
		})
	}
	fmt.Println("gee finished")

	r.Run(":9999")
}