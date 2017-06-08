package jr

import (
	"crypto/md5"
	"encoding/hex"
	// "reflect"
	"strconv"
)

type Session struct {
	ctx *Context
	M   map[string]interface{}
}

func newSession(ctx *Context) *Session {
	s := new(Session)
	s.ctx = ctx
	err := s.getPID()
	if err != nil {
		s.ctx.ELog(err)
		return nil
	}
	mc := s.ctx.Memcache()
	s.M = make(map[string]interface{})
	key := s.getKey()
	// s.ctx.Log.Print("SESSION KEY", key)
	mc.Get(key, &s.M)
	// s.ctx.Log.Print("mc get M", s.M)
	return s
}

func (s *Session) getPID() error {
	cookie := s.ctx.Cookie()
	var pid string
	pid, err := cookie.Get("PID")
	if err != nil {
		hasher := md5.New()
		hasher.Write([]byte(strconv.FormatInt(UniqID(), 10)))
		pid = hex.EncodeToString(hasher.Sum(nil))
		s.ctx.ILog("SESSION", pid)
		cookie.Set("PID", pid, s.ctx.config.SessionExpired)
	}
	s.ctx.Register["PID"] = pid
	return nil
}

func (s *Session) getKey() string {
	hasher := md5.New()
	hasher.Write([]byte(s.ctx.Register["PID"].(string)))
	hasher.Write([]byte(s.ctx.Request.UserAgent()))
	hasher.Write([]byte(s.ctx.RemoteAddr()))
	pid := hex.EncodeToString(hasher.Sum(nil))

	key := "SESSION_" + pid
	return key
}

func (s *Session) Commit() {
	// s.ctx.ILog("s", s.M)
	if len(s.M) > 0 {
		mc := s.ctx.Memcache()
		key := s.getKey()
		// s.ctx.ILog("k", key)
		mc.Set(key, s.M, 86400)
	}
}

func (s *Session) Get(name string) interface{} {
	if val, ok := s.M[name]; ok {
		// s.ctx.Log.Info("session get M:%v value:%v (%v)", name, val, reflect.TypeOf(val))
		return val
	}
	return nil
}

func (s *Session) Set(name string, value interface{}) {
	// c.ILog("name", name)
	// c.ILog("value", value)
	// s.ctx.Log.Info("session set M:%v value:%v (%v)", name, value, reflect.TypeOf(value))
	s.M[name] = value
}
