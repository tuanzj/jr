package jr

import (
	"strconv"
	"strings"
)

type (
	RouteFn func(c *Context) bool
	Route   struct {
		routeURL string
		method   string
		params   map[int]string
		routeFn  []RouteFn
		wildcard bool
	}
	Routes []Route
)

func initRoutes() Routes {
	return make(Routes, 0)
}

/*
func newRoute(routeURL string, routeFn RouteFn, method string) *Route {
    r := new(Route)
	s := strings.Split(routeURL, "/:")
	start := 0
	len_s := 0
	if strings.Index(s[0], "/") == 0 { // case / or /x or /c/a or /a/:b/:c
		r.routeURL = s[0]
		len_s = len(s) - 1
	} else { // case /:a/:b/:c ...
		len_s = len(s)
		r.routeURL =  "/"
		start = 1
	}
    r.params = make(map[int]string, 0)
	if len_s > 0 {
		for i, p := range s {
			if i - start < 0 { // prevent unpredictable extra params < 0
				continue
			}
			if p == "*" {
				r.wildcard = true
				break;
			}
			r.params[i - start] = p
		}
	}
	r.routeFn = routeFn
	r.method = method
    return r
}
*/
func newRoute(routeURL string, method string, routeFn ...RouteFn) *Route {
	r := new(Route)
	s := strings.Split(routeURL, "/:")
	start := 0

	if strings.Index(s[0], "/") == 0 {
		r.routeURL = s[0]
	} else {
		r.routeURL = "/"
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

/*
//stable
func newRoute(routeURL string, routeFn RouteFn, method string) *Route {
    r := new(Route)
	s := strings.Split(routeURL, "/:")
	start := 0
	if strings.Index(s[0], "/") == 0 {
		r.routeURL = s[0]
	} else {
		r.routeURL =  "/"
		start = 1
	}
    r.params = make(map[int]string, 0)
	if len(s) > 0 {
		for i, p := range s {
			if i - start < 0 { // prevent unpredictable extra params < 0
				continue
			}
			if p == "*" {
				r.wildcard = true
				break;
			}
			r.params[i - start] = p
		}
	}
	r.routeFn = routeFn
	r.method = method
    return r
}
*/
func (rts Routes) Match(c *Context) []RouteFn {
	url := c.Request.URL.RequestURI()
	for _, rt := range rts {
		//if i < 1 {
		//	continue
		//}
		if rt.method != "" && rt.method != "WS" && rt.method != c.Request.Method {
			continue
		}
		tmp := strings.Index(url, rt.routeURL)

		// c.ILog("url match" ,url)
		// c.ILog("routeURL" ,rt.routeURL)
		// c.ILog("tmp" ,tmp)
		// c.Log.Info("route %v %v %v", url, rt.routeURL, tmp)
		if tmp == 0 {
			tmp2 := strings.Split(url, "?") // exscind query part
			// c.ILog("tmp2" ,tmp2)
			// c.ILog("routeURL" ,rt.routeURL)
			if len(tmp2) > 0 {
				tmp3 := strings.Replace(tmp2[0], rt.routeURL, "", 1) // exscind match part
				if tmp3 != "" && strings.Index(tmp3, "/") != 0 {     // match with split url ex: /xa match with /x (except it full match tmp3 == "")
					if rt.routeURL != "/" || !rt.wildcard { // sua tam case /:* //TODO : kiem tra lai cho nay (day co the la mot ngoai le)
						continue
					}
				}
				s := strings.Split(tmp3, "/")

				//   	c.ILog("s" ,s)
				// c.ILog("ls" , len(s))
				// c.ILog("lp" , len(rt.params))
				// c.ILog("wc" , rt.wildcard)
				// c.ILog("par" , rt.params)
				if len(s) > 0 {
					if len(s) > len(rt.params) && !rt.wildcard { // prevent short url map with another longer
						continue
					}
					for i, p := range s {
						if p != "" {
							if _, ok := rt.params[i]; ok {
								c.Param[rt.params[i]] = p
							} else { //if rt.wildcard {
								c.Param[strconv.Itoa(i)] = p
							}
						}
					}
				}
			}
			return rt.routeFn
		}
	}
	return nil
	// return 404 default page
	// return rts[0].routeFn // return 404 default page
}

/*
func (rts Routes) Match(c *Context) RouteFn {
	url := c.Request.URL.RequestURI()
	for _, rt := range rts {
		//if i < 1 {
		//	continue
		//}
        if rt.method != "" && rt.method != c.Request.Method {
            continue
        }
		tmp := strings.Index(url, rt.routeURL) //
        // c.Log.Info("route %v %v %v", url, rt.routeURL, tm p)
		if tmp == 0 {
			tmp2 := strings.Split(url, "?") // exscind query parts
			if len(tmp2) > 0 {
            	c.ILog("tmp2" ,tmp2)
            	c.ILog("routeURL" ,rt.routeURL)
            	s := strings.Split(strings.Replace(tmp2[0], rt.routeURL, "", 1), "/") // exscind match part
            	c.ILog("s" ,s)
        		c.ILog("ls" , len(s))
        		c.ILog("lp" , len(rt.params))
        		c.ILog("wc" , rt.wildcard)
        		c.ILog("par" , rt.params)
            	len_s := len(s)
            	len_pr := len(rt.params)
	            if len_s == 1 {
	            	if s[0] == "" {
	            		if len_pr > 0 && !rt.wildcard {
	            			c.ILog("continue1", "")
	            			continue
	            		}
	            	} else {
	            		if len_s > len_pr { // prevent short url map with another longer
							if !rt.wildcard {
								c.ILog("continue2", "")
	            				continue
							}
						}
	            	}
	            } else if len_s > 1 {
	            	if len_s > len_pr { // prevent short url map with another longer
						if !rt.wildcard {
							c.ILog("continue3", "")
            				continue
						}
					}
	            }
	            /*
	            	if s[0] != "" { // || s[0] != " "
	            		 // else {
            				if len(s) > len(rt.params) && !rt.wildcard { // prevent short url map with another longer
	    						continue
	            			}
	            		// }
            		} else {
            			if len(rt.params) > 0 {
	            			c.ILog("continue", rt)
	            			continue
	            		}
            		}

	            *
	            	// c.Log.Info("len route %v %v %v", len(s) , len(rt.params), rt.params)
                for i, p := range s {
                    if p != "" {
                        if _, ok := rt.params[i]; ok {
                            c.Param[rt.params[i]] = p
                        } else {
                        	c.Param[strconv.Itoa(i)] = p
                        }
                    }
                }
	            	// c.ILog("match", rt)
	            // }
	        }
			return rt.routeFn
		}
	}
	 // return 404 default page
	return nil
	// return rts[0].routeFn // return 404 default page
}
*/
