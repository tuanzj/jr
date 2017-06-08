package jr

import (
	"gopkg.in/mgo.v2"
)

func newMGo(ctx *Context, host string) *mgo.Session {
	session, err := mgo.Dial(host)
    if err != nil {
            panic(err)
    }
    return session
}
