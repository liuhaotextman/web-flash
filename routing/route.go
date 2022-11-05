package routing

import (
	"net/http"
	httpFlash "web-flash/http"
)

type HandlerFunc func(*httpFlash.Context)

type Route struct {
	routeMap      map[string]HandlerFunc
	routeGroupMap map[string]*RouteGroup
}

func (route Route) RulePattern(method string, pattern string) string {
	return method + "-" + pattern
}

func (route *Route) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = route.RulePattern(method, pattern)
	route.routeMap[pattern] = handler
}

func New() *RouteGroup {
	route := &Route{
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
	pattern := route.RulePattern(ctx.Method, ctx.Path)
	handler, ok := route.routeMap[pattern]
	if !ok {
		ctx.Html(404, "not found")
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
