package gee

import "log"

// RouterGroup 路由组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
}

// Group 添加新的路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("GET", pattern, handlerFunc)
}

func (group *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("POST", pattern, handlerFunc)
}
