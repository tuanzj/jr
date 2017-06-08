package jr

type IURLFetch interface {
	Get(url string, query map[string]string, headers ...map[string]string) (string, error)
	Post(url string, data map[string]string, headers ...map[string]string) (string, error)
	PostJSON(urlString string, data interface{}, headers ...map[string]string) (string, error)
	PostXML(urlString string, data interface{}, headers ...map[string]string) (string, error)
	URLString(urlString string, parameters map[string]string) string
}

func newURLFetch(ctx *Context) IURLFetch {
	if ctx.isGAE() {
		return &GAEURLFetch{
			ctx: ctx,
		}
	}
	return &XURLFetch{
		ctx: ctx,
	}
}
