package jr

import (
    "google.golang.org/appengine/log"
)

type GAELog struct {
    ctx *Context
}

func (l *GAELog) Print(txt string, arg interface{}) {
    log.Infof(l.ctx.ctx, "JB ["+txt+"] : %v", arg)
}

func (l *GAELog) Info(txt string, args ...interface{}) {
    log.Infof(l.ctx.ctx, txt, args...)
}

func (l *GAELog) Error(txt string, args ...interface{}) {
    log.Errorf(l.ctx.ctx, txt, args...)
}

func (l *GAELog) Warning(txt string, args ...interface{}) {
    log.Warningf(l.ctx.ctx, txt, args...)
}

func (l *GAELog) Debug(txt string, args ...interface{}) {
    log.Debugf(l.ctx.ctx, txt, args...)
}

func (l *GAELog) Critical(txt string, args ...interface{}) {
    log.Criticalf(l.ctx.ctx, txt, args...)
}
