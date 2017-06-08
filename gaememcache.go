package jr

import (
	// "encoding/json"
	"google.golang.org/appengine/memcache"
	// "encoding"
	"bytes"
	"encoding/gob"
	"time"
)

type GAEMemcache struct {
	ctx *Context
}

func (m *GAEMemcache) Set(name string, value interface{}, expire int) error {
	// v, err := json.Marshal(value)
	// if err != nil {
	//     return err
	// }

	buffer := bytes.Buffer{}
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        name,
		Value:      buffer.Bytes(),
		Expiration: time.Duration(expire) * time.Second,
	}
	m.ctx.Log.Info("memcache [%v] added", item.Key)
	memcache.Set(m.ctx.ctx, item)
	return nil
}

func (m *GAEMemcache) Get(name string, rs interface{}) error {
	item, err := memcache.Get(m.ctx.ctx, name)
	if err == memcache.ErrCacheMiss {
		return err
		m.ctx.ELog(err)
	} else if err != nil {
		m.ctx.ELog(err)
		return err
	} else {
		// err = json.Unmarshal(item.Value, rs)
		// if err != nil {
		//     return err
		// }

		buffer := bytes.Buffer{}
		buffer.Write(item.Value)
		data := gob.NewDecoder(&buffer)
		err = data.Decode(rs)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func ToGOB64(m SX) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) SX {
	m := SX{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}
*/
