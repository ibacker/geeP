package gee

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	// 继承RouterGroup 所有属性和方法
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	// 添加顶层路由分组
	engine.groups = []*RouterGroup{engine.RouterGroup}
	fmt.Println("new Engine: ", engine)
	return engine
}

func (engine *Engine) addRoute(method string, pattern string, handlerFunc HandlerFunc) {
	engine.router.addRoute(method, pattern, handlerFunc)
}

func (engine *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	engine.addRoute("GET", pattern, handlerFunc)
}

func (engine *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	engine.addRoute("POST", pattern, handlerFunc)
}

// 实现了 handler 方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	engine.router.handler(ctx)
}

func (engine *Engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, engine))
}