package jr

type IMemcache interface {
	Set(name string, value interface{}, expire int) error
	Get(name string, rs interface{}) error
}

func newMemcache(ctx *Context) IMemcache {
	if ctx.isGAE() {
		return &GAEMemcache{
			ctx: ctx,
		}
	}
	return NewXMemcache(ctx)
}
