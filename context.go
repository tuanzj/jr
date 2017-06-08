package jr

import (
	"encoding/json"
	"encoding/xml"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"html/template"
	"net"
	"net/http"
	// "net/url"
	"strings"
)

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	Param      map[string]string
	Register   map[string]interface{}
	remoteAddr string
	ctx        context.Context
	datastore  []IDB
	memcache   IMemcache
	session    *Session
	cookie     *Cookie
	Log        ILog
	config     *Config
}

func NewContext(config *Config, w http.ResponseWriter, r *http.Request) *Context {
	c := new(Context)
	c.Param = make(map[string]string)
	c.Register = make(map[string]interface{})
	c.Response = w
	c.Request = r
	c.config = config
	if c.isGAE() {
		c.ctx = appengine.NewContext(r)
	}
	c.Log = newLog(c)
	c.datastore = make([]IDB, 0)
	return c
}

// func (c *Context) getContext() Context {
// 	return c
// }

// func (c *Context) getGAEContext() context.Context {
// 	return c.ctx
// }

func (c *Context) isGAE() bool {
	return c.config.GAE
}

func (c *Context) Finalize() {
	if c.session != nil {
		c.session.Commit()
	}
	n := len(c.datastore)
	if n > 0 {
		for i := 0; i < n; i++ {
			c.datastore[i].Disconnect()
		}
	}
}

func (c *Context) DBDriver(connect ...string) IDB {
	db := newDB(c, connect...)
	c.datastore = append(c.datastore, db)
	return db
}

func (c *Context) Memcache() IMemcache {
	if c.memcache == nil {
		c.memcache = newMemcache(c)
	}
	return c.memcache
}

func (c *Context) URLFetch() IURLFetch {
	return newURLFetch(c)
}

func (c *Context) Session() *Session {
	if c.session == nil {
		c.session = newSession(c)
	}
	return c.session
}
func (c *Context) Cookie() *Cookie {
	if c.cookie == nil {
		c.cookie = newCookie(c)
	}
	return c.cookie
}
func (c *Context) ILog(tag string, arg interface{}) {
	c.Log.Print(tag, arg)
}

func (c *Context) ELog(arg interface{}) {
	c.Log.Error("", arg)
}

// func (c *Context) MGo(host string) *mgo.Session {
// 	return newMGo(c, host)
// }

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) SetHeader(name string, content string) {
	c.Response.Header().Set(name, content)
}
func (c *Context) CORS(domain ...string) {
	if origin := c.Request.Header.Get("Origin"); origin != "" {
		c.SetHeader("Access-Control-Allow-Origin", origin)
	}
	c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.SetHeader("Access-Control-Allow-Headers",
		"Accept, X-Requested-With, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.SetHeader("Access-Control-Allow-Credentials", "true")
	// if origin := c.Request.Header.Get("Origin"); origin != "" {
	// 	c.SetHeader("Access-Control-Allow-Origin", origin)
	// } else {
	// }
	/*
		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		// c.SetHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
		c.SetHeader("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Credentials")
	*/
	// }
	// if c.Request.Method == "OPTIONS" {
	// 	return
	// }
}

func (c *Context) Query(name string) string {
	return c.Request.FormValue(name)
}
func (c *Context) QueryList(name string) []string {
	v := c.Request.URL.Query()
	return v[name]
}

// func (c *Context) QueryList() url.Values {
// 	u, _ := url.Parse(c.Request.URL.String())
// 	queryParams := u.Query()
// 	return queryParams
// }

func (c *Context) RequestJSONData(data interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}
func (c *Context) RequestXMLData(data interface{}) error {
	decoder := xml.NewDecoder(c.Request.Body)
	// var data map[string]interface{}
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func (c *Context) GetVar(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) PostVar(name string) string {
	return c.Request.FormValue(name)
}

func (c *Context) RemoteAddr() string {
	if c.remoteAddr == "" {
		for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
			addresses := strings.Split(c.Request.Header.Get(h), ",")
			// march from right to left until we get a public address
			// that will be the address right before our proxy.
			// c.ILog("ip", addresses)
			for i := len(addresses) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(addresses[i])
				// header can contain spaces too, strip those out.
				realIP := net.ParseIP(ip)
				if !realIP.IsGlobalUnicast() { //} || isPrivateSubnet(realIP) {
					// bad address, go to next
					continue
				}
				c.remoteAddr = ip
			}
		}
		c.remoteAddr = c.Request.RemoteAddr
	}
	return c.remoteAddr
}

func (c *Context) JSON(data interface{}) bool {
	c.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(c.Response).Encode(data); err != nil {
		c.ELog(err)
		return c.Abort(500)
	}
	return true
}

func (c *Context) XML(data interface{}) bool {
	c.Response.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if err := xml.NewEncoder(c.Response).Encode(data); err != nil {
		c.ELog(err)
		return c.Abort(500)
	}
	return true
}

func (c *Context) IsXMLHttpRequest() bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
		c.Request.Header.Get("Origin") != ""
}

func (c *Context) Next() bool {
	return true
}

func (c *Context) Abort(code int) bool {
	c.Response.WriteHeader(code)
	if !c.IsXMLHttpRequest() {
		if _, ok := c.config.HTTPError[code]; ok {
			return c.config.HTTPError[code](c)
		}
	}
	return false
}

func (c *Context) Text(str string) bool {
	if _, err := c.Response.Write([]byte(str)); err != nil {
		c.ELog(err)
	}
	return true
}

func (c *Context) Redirect(url string) bool {
	http.Redirect(c.Response, c.Request, url, http.StatusFound)
	return true
}

func (c *Context) Render(name string, tpl string, args ...interface{}) bool {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}
	template := template.Must(template.New(name).Parse(tpl))
	err := template.Execute(c.Response, arg)
	if err != nil {
		c.Log.Error("[TEMPLATE] %v", err)
	}
	return true
}

func (c *Context) RenderFile(filename string, args ...interface{}) bool {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}
	template := template.Must(template.ParseFiles(filename))
	c.ILog("TPL", template)
	err := template.Execute(c.Response, arg)
	if err != nil {
		c.ELog(err)
	}
	return true
}

func (c *Context) Push(uri string) error {
	// if pusher, ok := c.Response.(http.Pusher); ok {
	// 	if err := pusher.Push(uri, nil); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

/*

type ipRange struct {
    start net.IP
    end net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
    // strcmp type byte comparison
    if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) <= 0 {
        return true
    }
    return false
}

var privateRanges = []ipRange{
    ipRange{
        start: net.ParseIP("10.0.0.0"),
        end:   net.ParseIP("10.255.255.255"),
    },
    ipRange{
        start: net.ParseIP("100.64.0.0"),
        end:   net.ParseIP("100.127.255.255"),
    },
    ipRange{
        start: net.ParseIP("172.16.0.0"),
        end:   net.ParseIP("172.31.255.255"),
    },
    ipRange{
        start: net.ParseIP("192.0.0.0"),
        end:   net.ParseIP("192.0.0.255"),
    },
    ipRange{
        start: net.ParseIP("192.168.0.0"),
        end:   net.ParseIP("192.168.255.255"),
    },
    ipRange{
        start: net.ParseIP("198.18.0.0"),
        end:   net.ParseIP("198.19.255.255"),
    },
}


// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
    // my use case is only concerned with ipv4 atm
    if ipCheck := ipAddress.To4(); ipCheck != nil {
        // iterate over all our ranges
        for _, r := range privateRanges {
            // check if this ip is in a private range
            if inRange(r, ipAddress){
                return true
            }
        }
    }
    return false
}
*/
