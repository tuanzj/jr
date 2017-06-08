package jr

type ILog interface {
    Print(txt string, arg interface{})
    Info(txt string, args ...interface{})
    Error(txt string, args ...interface{})
    Warning(txt string, args ...interface{})
    Debug(txt string, args ...interface{})
    Critical(txt string, args ...interface{})
}

func newLog(ctx *Context) ILog {
    if ctx.isGAE() {
        return &GAELog {
            ctx: ctx,
        }
    }
    return &XLog {
        ctx: ctx,
    }
}
