package jr

type IAMPQ interface {
}

func newAMPQ(ctx *Context, connect string) IAMPQ {
	o := &XAMPQ{
		ctx: ctx,
	}
	// o.Connect(connect)
	return o
}
