package jr

type Config struct {
	GAE              bool
	DevENV           bool
	LogFile          string
	MemcacheHost     string
	SessionExpired   int
	DefaultCacheTime int
	HTTPError        map[int]RouteFn
}

func NewConfig() *Config {
	config := &Config{
		GAE:              false,
		DevENV:           false,
		LogFile:          "",
		MemcacheHost:     "",
		SessionExpired:   3600,
		DefaultCacheTime: 300,
		HTTPError:        make(map[int]RouteFn),
	}
	config.HTTPError[404] = http404Error
	config.HTTPError[500] = http500Error

	return config
}

func http404Error(c *Context) bool {
	return c.Text(`<!doctype html>
						<html>
						<head>
							<title>PAGE NOT FOUND</title>
						</head>
						<body>
							<h1>PAGE NOT FOUND</h1>
							<b>This might be because:</b>
							<ul>
								<li>You have typed the web address incorrectly, or</li>
								<li>the page you were looking for may have been moved, updated or deleted.</li>
							</ul>
						</body>
						</html>`)
}
func http500Error(c *Context) bool {
	return c.Text(`<!doctype html>
						<html>
						<head>
							<title>INTERNAL SERVER ERROR</title>
						</head>
						<body>
							<h1>INTERNAL SERVER ERROR</h1>
						</body>
						</html>`)
}

/*
config.HTTPError[404] = `<!doctype html>
						<html>
						<head>
							<title>PAGE NOT FOUND</title>
						</head>
						<body>
							<h1>PAGE NOT FOUND</h1>
							<b>This might be because:</b>
							<ul>
								<li>You have typed the web address incorrectly, or</li>
								<li>the page you were looking for may have been moved, updated or deleted.</li>
							</ul>
						</body>
						</html>`
	config.HTTPError[500] = `<!doctype html>
						<html>
						<head>
							<title>INTERNAL SERVER ERROR</title>
						</head>
						<body>
							<h1>INTERNAL SERVER ERROR</h1>
						</body>
						</html>`

*/
