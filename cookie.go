package jr

import "net/http"

type Cookie struct {
	ctx *Context
}

func newCookie(ctx *Context) *Cookie {
	c := new(Cookie)
	c.ctx = ctx
	return c
}

func (c *Cookie) Get(name string) (string, error) {
	cookie, err := c.ctx.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
func (c *Cookie) Set(name string, value string, expire int) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = "/"
	cookie.MaxAge = expire
	http.SetCookie(c.ctx.Response, cookie)
}

/*
func (c *Cookie) Get(name string, result interface{}) error {
	cookie, err := c.ctx.Request.Cookie(name)
	if err != nil {
		// c.ctx.Log.Error("COOKIE %v", err)
		return err
	}
	v, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), result)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cookie) Set(name string, value interface{}, expire int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = name
	// c.ctx.ILog("raw cookie value", value)
	// c.ctx.ILog("cookie value", string(v[:len(v)]))
	// cach convert bytes[] to string
	cookie.Value = url.QueryEscape(string(v[:len(v)])) //ByteToString(v)
	cookie.Path = "/"
	cookie.MaxAge = expire
	http.SetCookie(c.ctx.Response, cookie)
	return nil
}
*/
