package jr

import (
	"bytes"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
)

type XMemcache struct {
	ctx *Context
	mc  *memcache.Client
}

func NewXMemcache(ctx *Context) IMemcache {
	mc := &XMemcache{
		ctx: ctx,
		mc:  memcache.New("127.0.0.1:11211"),
	}
	return mc
}
func (m *XMemcache) isUseMemcache() bool {
	return m.ctx.config.MemcacheHost != ""
}

func (m *XMemcache) Set(name string, value interface{}, expire int) error {
	if m.isUseMemcache() {

		buffer := bytes.Buffer{}
		encode := gob.NewEncoder(&buffer)
		err := encode.Encode(value)
		if err != nil {
			return err
		}

		// v, err := json.Marshal(value)
		// if err != nil {
		// 	return err
		// }

		m.mc.Set(&memcache.Item{Key: name, Value: buffer.Bytes(), Expiration: int32(expire)}) // []byte(v)
	}
	return nil
}

func (m *XMemcache) Get(name string, rs interface{}) error {
	if m.isUseMemcache() {
		item, err := m.mc.Get(name)
		if err != nil {
			return err
			m.ctx.Log.Error("MEMCACHE: item [%v] error - %v", name, err)
		} else {
			buffer := bytes.Buffer{}
			buffer.Write(item.Value)
			data := gob.NewDecoder(&buffer)
			err = data.Decode(rs)
			if err != nil {
				return err
			}
			// json.Unmarshal(item.Value, rs)
		}
	}
	return nil
}
