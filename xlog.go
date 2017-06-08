package jr

import (
    "log"
    "os"
)

type XLog struct {
    ctx *Context
}

func (l *XLog) XLogPrint(txt string, args ...interface{}) {
    if l.ctx.config.LogFile != "" {
        f, err := os.OpenFile(l.ctx.config.LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        if err != nil {
        }
        defer f.Close()
        log.SetOutput(f)
    }
    log.Printf(txt, args...)
}
func (l *XLog) Print(txt string, arg interface{}) {
    l.XLogPrint("JB ["+txt+"] : %v", arg)
}

func (l *XLog) Info(txt string, args ...interface{}) {
    l.XLogPrint("[INFO]" + txt, args...)
}

func (l *XLog) Error(txt string, args ...interface{}) {
    l.XLogPrint("[ERROR]" + txt, args...)
}

func (l *XLog) Warning(txt string, args ...interface{}) {
    l.XLogPrint("[WARNING]" + txt, args...)
}

func (l *XLog) Debug(txt string, args ...interface{}) {
    l.XLogPrint("[DEBUG]" + txt, args...)
}

func (l * XLog) Critical(txt string, args ...interface{}) {
    l.XLogPrint("[CRITICAL]" + txt, args...)
}
