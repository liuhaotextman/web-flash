package routing

import "net/http"

type RouteGroup struct {
	prefix string
	parent *RouteGroup
	middlewares []HandlerFunc
	route *Route
}

func (group *RouteGroup) addRoute(method string, pattern string, handler HandlerFunc){
	pattern = group.prefix + pattern
	group.route.addRoute(method, pattern, handler)
	group.route.AddRouteGroupMap(method, pattern, group)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouteGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}

func (group *RouteGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

func (group *RouteGroup) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		middlewares: group.middlewares,
		route: group.route,
	}
}

func (group *RouteGroup) Use(handler HandlerFunc) {
	group.middlewares = append(group.middlewares, handler)
}

func (group *RouteGroup) Run(address string) {
	http.ListenAndServe(address, group.route)
}