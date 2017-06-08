package jr

import (
	"net/http"
)

// type MiddlewareFn func(c *Context) bool

type Handler struct {
	routes     Routes
	middleware []RouteFn
	context    *Context
	config     *Config
}

func New(config ...*Config) *Handler { //time string
	hd := new(Handler)
	hd.routes = initRoutes()
	hd.middleware = make([]RouteFn, 0)

	if len(config) > 0 {
		hd.config = config[0]
	} else {
		// defaut config
		hd.config = NewConfig()
	}

	// default route 404 error page
	// hd.addRoute("", httpError404, "")

	return hd
}

// var context *Context

// func CTX() *Context {
// 	return context
// }

func (hd *Handler) Run(addr ...string) {
	if hd.config.GAE {
		http.Handle("/", hd)
		return
	}
	address := resolveAddress(addr)
	http.ListenAndServe(address, hd)
}

// http entry func
func (hd *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// each request has one context
	context := NewContext(hd.config, w, r)
	routefn := hd.routes.Match(context)
	if routefn == nil {
		// hd.config.httpError404(context)
		context.Abort(404)
		return
	}
	// run middleware func one by one with context
	for _, mw := range hd.middleware {
		rs := mw(context)
		if !rs {
			return
		}
	}
	// context.ILog("routefn", routefn)
	for _, rt := range routefn {
		rs := rt(context)
		if !rs {
			return
		}
	}
	context.Finalize()
}

// add middleware
func (hd *Handler) Use(fn RouteFn) {
	hd.middleware = append(hd.middleware, fn)
}

func (hd *Handler) addRoute(routeURL string, method string, routeFn ...RouteFn) {
	rt := newRoute(routeURL, method, routeFn...)
	hd.routes = append(hd.routes, *rt)
}

func (hd *Handler) Route(url string, fn ...RouteFn) {
	hd.addRoute(url, "", fn...)
}
func (hd *Handler) GET(url string, fn ...RouteFn) {
	hd.addRoute(url, "GET", fn...)
}
func (hd *Handler) POST(url string, fn ...RouteFn) {
	hd.addRoute(url, "POST", fn...)
}
func (hd *Handler) PUT(url string, fn ...RouteFn) {
	hd.addRoute(url, "PUT", fn...)
}
func (hd *Handler) DELETE(url string, fn ...RouteFn) {
	hd.addRoute(url, "DELETE", fn...)
}
