package gee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type router struct {
	// 请求路径根节点, 根据 请求类型区分
	roots map[string]*node
	// 请求响应方法
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc),
		roots: make(map[string]*node)}
}

// 解析请求路径
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 单一路径仅允许一个通配符*
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	_, ok := r.roots[method]

	// 首次新增
	if !ok {
		fmt.Println("addRoute-First Add Method", method)
		fmt.Println("addRoute-First Add pattern", pattern)
		r.roots[method] = &node{}
	}
	// 插入路由节点信息
	r.roots[method].insert(pattern, parts, 0)
	fmt.Println("addRoute-insert ", r.roots[method])
	key := method + "_" + pattern
	// 插入处理方法
	r.handlers[key] = handler
	fmt.Println("addRoute finish, roots: ", r.roots[method])
	fmt.Println("addRoute finish, handlers: ", r.handlers)
}

// 查询路由节点
// 路由参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {

	fmt.Println("getRoute-path", path)
	searchParts := parsePattern(path)
	params := make(map[string]string)
	// 查询请求类型是否存在
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	fmt.Println("getRoute-searchParts", searchParts)

	// 查询路由节点
	n := root.search(searchParts, 0)
	fmt.Println("getRoute-root.root", root)
	fmt.Println("getRoute-root.search", n)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 路由参数
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 通配符 将*后面的字符认为参数
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 路由处理
//
//	func (r *router) handle(c *Context) {
//		n, params := r.getRoute(c.Method, c.Path)
//		if n != nil {
//			c.Params = params
//			// 通过节点路径匹配
//			key := c.Method + "_" + n.pattern
//			// 找到对应处理方法并执行
//			fmt.Println("handler", key)
//			r.handlers[key](c)
//		} else {
//			c.String(http.StatusNotFound, "404 NOT FOUND : %s\n", c.Path)
//		}
//		// 执行中间件方法
//		c.Next()
//	}
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "_" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
