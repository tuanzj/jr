package jr

/*
import "strings"

type (
	MiddlewareFn func(c *Context) bool
	Middleware   struct {
		path     string
		fn       MiddlewareFn
		wildcard bool
	}
	Middlewares []Middleware
)

func initMiddlewares() Middlewares {
	return make(Middlewares, 0)
}

func newMiddleware(path string, fn MiddlewareFn) *Middleware {
	// r := new(Middleware)
	// s := strings.Split(path, "/:")
	// start := 0

	// if strings.Index(s[0], "/") == 0 {
	// r.routeURL = s[0]
	// } else {
	// r.routeURL = "/"
	// start = 1
	// }
	// r.params = make(map[int]string, 0)
	// if len(s) > 0 {
	// for i, p := range s {
	// if i-start < 0 { // prevent unpredictable extra params < 0
	// continue
	// }
	// if p == "*" {
	// r.wildcard = true
	// break
	// }
	// r.params[i-start] = p
	// }
	// }
	// r.path = path
	// r.fn = routeFn
	// r.method = method
	// return r

	r := new(MiddlewareFn)
	s := strings.Split(path, "/:")
	start := 0

	if strings.Index(s[0], "/") == 0 {
		r.path = s[0]
	} else {
		r.path = "/"
		start = 1
	}
	r.params = make(map[int]string, 0)
	if len(s) > 0 {
		for i, p := range s {
			if i-start < 0 { // prevent unpredictable extra params < 0
				continue
			}
			if p == "*" {
				r.wildcard = true
				break
			}
			r.params[i-start] = p
		}
	}
	r.routeFn = routeFn
	r.method = method
	return r
}
func (rts Middlewares) Match(c *Context) MiddlewareFn {
	url := c.Request.URL.RequestURI()
	for _, rt := range rts {
		tmp := strings.Index(url, rt.routeURL)
		if tmp == 0 {
			return rt.fn
		}
	}
	return nil
}
*/
