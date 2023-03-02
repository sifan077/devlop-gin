package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	// 根节点
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	// 将请求的url分为字符串数组
	vs := strings.Split(pattern, "/")
	// 声明一个匹配路由字符串切片
	parts := make([]string, 0)
	// 遍历url分解的字符串数组,如果不为空字符串则加入待匹配数组,如果是通配符则直接跳出循环
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	// 分解url路径字符串
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	// 从根节点获取节点
	_, ok := r.roots[method]
	// 如果不存在节点,则添加一个节点
	if !ok {
		r.roots[method] = &node{}
	}
	// 在根节点中插入当前代匹配的字符串和分解的字符串数组
	r.roots[method].insert(pattern, parts, 0)
	// 生命key对用的处理方法
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 分解url路径字符串
	searchParts := parsePattern(path)
	// 声明参数切片
	params := make(map[string]string)
	// 从根节点获取节点
	root, ok := r.roots[method]
	// 如果不存在则返回null
	if !ok {
		return nil, nil
	}
	// 在根节点搜索需要搜索的字符串数组
	n := root.search(searchParts, 0)
	// 如果搜索到了节点
	if n != nil {
		// 分解寻找到的节点的代匹配路由字符串
		parts := parsePattern(n.pattern)
		// 遍历代匹配的路由字符串
		for index, part := range parts {
			// 如果是 ':',则向params的map对应的key添加value
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
