package routing

import (
	"net/http"
	"strings"
	httpFlash "web-flash/http"
)

type HandlerFunc func(*httpFlash.Context)

type Route struct {
	roots         map[string]*node
	routeMap      map[string]HandlerFunc
	routeGroupMap map[string]*RouteGroup
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	patterns := make([]string, 0)
	for _, part := range vs {
		if part != "" {
			patterns = append(patterns, part)
		}
	}

	return patterns
}

func (route Route) RulePattern(method string, pattern string) string {
	return method + "-" + pattern
}

func (route *Route) addRoute(method string, pattern string, handler HandlerFunc) {
	_, ok := route.roots[method]
	if !ok {
		route.roots[method] = &node{}
	}
	parts := parsePattern(pattern)
	route.roots[method].insert(pattern, parts, 0)
	pattern = route.RulePattern(method, pattern)
	route.routeMap[pattern] = handler
}

func (route *Route) getRoute(method string, pattern string) (*node, map[string]string) {
	root, ok := route.roots[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	searchParts := parsePattern(pattern)
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
		}
		return n, params
	}

	return nil, nil
}

func New() *RouteGroup {
	route := &Route{
		roots:         make(map[string]*node),
		routeMap:      make(map[string]HandlerFunc),
		routeGroupMap: make(map[string]*RouteGroup),
	}
	return &RouteGroup{
		prefix: "",
		parent: nil,
		route:  route,
	}
}

func (route *Route) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := httpFlash.NewContext(writer, request)
	route.Handle(context)
}

func (route *Route) Handle(ctx *httpFlash.Context) {
	n, params := route.getRoute(ctx.Method, ctx.Path)
	if n == nil {
		ctx.Html(404, "route not found")
		return
	}
	ctx.Params = params
	pattern := n.pattern
	pattern = route.RulePattern(ctx.Method, pattern)
	handler, ok := route.routeMap[pattern]
	if !ok {
		ctx.Html(404, "handler not found")
		return
	}
	//handler(ctx)
	handlers := make([]HandlerFunc, 0)
	handlers = append(handlers, handler)
	if routeGroup, ok := route.routeGroupMap[pattern]; ok {
		for i := 0; i < len(routeGroup.middlewares); i++ {
			handlers = append(handlers, routeGroup.middlewares[i])
		}
	}
	s := len(handlers)
	for i := 0; i < s; i++ {
		handlers[i](ctx)
	}
}

func (route *Route) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix: prefix,
		parent: nil,
		route:  route,
	}
}

func (route *Route) AddRouteGroupMap(method string, pattern string, group *RouteGroup) {
	pattern = route.RulePattern(method, pattern)
	route.routeGroupMap[pattern] = group
}
