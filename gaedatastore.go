package jr

import (
	// "google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"strconv"
	// "reflect"
	// "golang.org/x/net/context"
)

// type GAEDO struct {
// 	key datastore.Key
// }

type IGAEDO interface {
	SetID(id interface{})
	GetID() interface{}
}

type IGAEModel interface {
	Load(p []datastore.Property) error
	Save() ([]datastore.Property, error)
}

// type GAEDO struct {
// 	key datastore.Key
// }

// func (d *GAEDO) SetID(id interface{}) {
// 	d.id = id
// }

// func (d *GAEDO) GetID() string {
// 	return d.id
// }

/*
Note ko xoa
Tai sao DO lai la &data con interface thi khong co vi interface la object
*/
type GAEDatastore struct {
	c *Context
}

func (d *GAEDatastore) Disconnect() error {
	// d.c.ILog(" datastore", "disconnect")
	return nil
}
func (d *GAEDatastore) Connect(connect ...string) error {
	// d.c.ILog(" datastore", "connect")
	return nil
}

// var CX *Context
/*
func (d *GAEDatastore) NewKey(kind, id string) IKey { //, parent IKey
    // var parent *datastore.Key
    var key *datastore.Key
    if intID, err := strconv.Atoi(id); err == nil {
        key = datastore.NewKey(d.c.ctx, kind, "", int64(intID), nil)
    } else {
        key = datastore.NewKey(d.c.ctx, kind, id, 0, nil)
    }
    return key
}
func (d *GAEDatastore) NewIncompleteKey(kind string) IKey {
    key := datastore.NewIncompleteKey(d.c.ctx, kind, nil)
    return key
}

func (d *GAEDatastore) DecodeKey(encoded string) (IKey, error) {
    return datastore.DecodeKey(encoded)
}
*/
// func (d *GAEDatastore) NewQuery(name string) map[string]interface{} {
//     return datastore.NewQuery(name)
// }

func (d *GAEDatastore) getKey(collection string, id ...interface{}) *datastore.Key { //, parent IKey
	var key *datastore.Key
	if len(id) == 1 {
		if idInt64, ok := id[0].(int64); ok {
			key = datastore.NewKey(d.c.ctx, collection, "", idInt64, nil)
		} else if idInt, ok := id[0].(int); ok {
			key = datastore.NewKey(d.c.ctx, collection, "", int64(idInt), nil)
		} else if idString, ok := id[0].(string); ok {
			key = datastore.NewKey(d.c.ctx, collection, idString, 0, nil)
		}
	} else {
		key = datastore.NewIncompleteKey(d.c.ctx, collection, nil)
	}
	return key
}

func (d *GAEDatastore) getQuery(collection string, selector map[string]interface{}) *datastore.Query {
	q := datastore.NewQuery(collection)
	useCursor := false
	var cursorName string
	var cursorPage int
	mc := d.c.Memcache()
	for k, v := range selector {
		switch k {
		case "$offset":
			if offset, ok := v.(int); ok {
				q = q.Offset(offset)
			}
		case "$limit":
			if limit, ok := v.(int); ok {
				q = q.Limit(limit)
			}
		case "$sort":
			if order, ok := v.(string); ok {
				q = q.Order(order)
			} else if orders, ok := v.([]string); ok {
				for i := 0; i < len(orders); i++ {
					q = q.Order(orders[i])
				}
			}
		case "$project":
			if projects, ok := v.([]string); ok {
				q = q.Project(projects...)
			}
		case "$keysonly":
			q = q.KeysOnly()
		case "$distinct":
			q = q.Distinct()
		case "$cursor":
			cursorPattern, ok := v.([]interface{})
			if !ok {
				break
			}
			// cursorName, ok := cursorPattern[0].(string) dau : o day thi se tao mot bien moi local, theng tao ngon ngu ngu lo
			cursorName, ok = cursorPattern[0].(string)
			// d.c.ILog("cursorName0", cursorName)
			if !ok {
				break
			}
			cursorPage, ok = cursorPattern[1].(int)
			// d.c.ILog("cursorPage0", cursorPage)
			if !ok {
				break
			}
			useCursor = true
			cursorKey := cursorName + "--" + strconv.Itoa(cursorPage)
			// d.c.ILog("cursorKey0", cursorKey)
			cursorString := ""
			err := mc.Get(cursorKey, &cursorString)
			if err != nil {
				cursorPage = 1
				d.c.ELog(err)
				break
			}
			cursor, err := datastore.DecodeCursor(cursorString)
			if err != nil {
				cursorPage = 1
				d.c.ELog(err)
				break
			}
			q = q.Start(cursor)
		default:
			q = q.Filter(k, v)
		}
	}
	if useCursor {
		cursorPage++
		cursorKey := cursorName + "--" + strconv.Itoa(cursorPage)
		cursorString := ""
		err := mc.Get(cursorKey, &cursorString)
		if err != nil {
			t := q.Run(d.c.ctx)
			for {
				_, err := t.Next(nil)
				if err == datastore.Done {
					break
				}
			}
			if cursor, err := t.Cursor(); err == nil {
				err = mc.Set(cursorKey, cursor.String(), d.c.config.DefaultCacheTime)
				if err != nil {
					d.c.ELog(err)
				}
			}
		}
	}
	return q
}

func (d *GAEDatastore) Find(collection string, id interface{}, dst interface{}) error {
	key := d.getKey(collection, id)
	if err := datastore.Get(d.c.ctx, key, dst); err != nil {
		d.c.ELog(err)
		return err
	}
	return nil
}

func (d *GAEDatastore) FindAll(collection string, selector map[string]interface{}, dst interface{}) error {
	q := d.getQuery(collection, selector)
	_, err := q.GetAll(d.c.ctx, dst)
	if err != nil {
		return err
	}
	return nil
}

func (d *GAEDatastore) InsertAll(collection string, docs ...interface{}) error {
	for i := 0; i < len(docs); i++ {
		d.Insert(collection, docs[i])
	}
	return nil
}

func (d *GAEDatastore) keyId(key *datastore.Key) interface{} {
	id := key.IntID()
	if id != 0 {
		return id
	}
	return key.StringID()
}

func (d *GAEDatastore) Insert(collection string, doc interface{}, id ...interface{}) error {
	key := d.getKey(collection, id...)
	if len(id) == 1 {
		doc.(IGAEDO).SetID(d.keyId(key))
	}
	key, err := datastore.Put(d.c.ctx, key, doc)
	if err != nil {
		d.c.ELog(err)
		return err
	}
	if len(id) == 0 {
		doc.(IGAEDO).SetID(d.keyId(key))
		_, err = datastore.Put(d.c.ctx, key, doc)
		if err != nil {
			d.c.ELog(err)
			return err
		}
	}
	return nil
}

func (d *GAEDatastore) Update(collection string, doc interface{}, id interface{}) error {
	key := d.getKey(collection, id)
	key, err := datastore.Put(d.c.ctx, key, doc)
	if err != nil {
		d.c.ELog(err)
		return err
	}
	return nil
}

func (d *GAEDatastore) UpdateAll(collection string, selector map[string]interface{}, doc interface{}) error {
	// _, err = datastore.PutMulti(ctx, []*datastore.Key{k1, k2, k3}, []interface{}{e1, e2, e3})
	return nil
}

func (d *GAEDatastore) Upsert(collection string, doc interface{}) error {
	id := doc.(IGAEDO).GetID()
	var err error
	if id == "" {
		err = d.Insert(collection, doc)
	} else {
		err = d.Update(collection, doc, id)
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *GAEDatastore) Remove(collection string, id interface{}) error {
	key := d.getKey(collection, id)
	d.c.ILog("key", key)
	err := datastore.Delete(d.c.ctx, key)
	if err != nil {
		d.c.ELog(err)
		return err
	}
	return nil
}

func (d *GAEDatastore) RemoveAll(collection string, selector map[string]interface{}) error {
	return nil
}

/*
func (d *GAEDatastore) FindByKey(key IKey, id string) (DO) {
    key := new(datastore.Key)
    rs := make(DO) // DO is not a struct, it is a map
    rs["$id"] = id
    rs["$key"] = key.Encode()
    if intID, err := strconv.Atoi(id); err == nil {
        key = datastore.NewKey(d.c.ctx, name, "", int64(intID), nil)
    } else {
        key = datastore.NewKey(d.c.ctx, name, id, 0, nil)
    }
    if err := datastore.Get(d.c.ctx, key, &rs); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    return rs
}

func (d *GAEDatastore) GetAll(q *datastore.Query) ([]DO) {
    rs := make([]DO, 0)
    keys, err := q.GetAll(d.c.ctx, &rs)
    if err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    // d.c.Log.Print("rs", rs)
    for i := 0; i < len(rs); i++ {
        rs[i]["$key"] = keys[i].Encode()
        rs[i]["$id"] = keys[i].IntID()
        if rs[i]["$id"].(int64) == 0 {
            rs[i]["$id"] = keys[i].StringID()
        // } else {
            // rs[i]["_id"] = strconv.FormatInt(rs[i]["_id"], 10)
        }
    }
    return rs
}

func (d *GAEDatastore) GetByKey(name string, keystring string) (DO) {
    // key := new(datastore.Key)
    rs := make(DO) // DO is not a struct, it is a map
    // rs["$id"] = id
    key, err := d.DecodeKey(keystring)
    if err != nil {
        return rs
    }
    // if intID, err := strconv.Atoi(id); err == nil {
        // key = datastore.NewKey(d.c.ctx, name, "", int64(intID), nil)
    // } else {
        // key = datastore.NewKey(d.c.ctx, name, id, 0, nil)
    // }
    if err := datastore.Get(d.c.ctx, key, &rs); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    } else {
        rs["$key"] = keystring
    }
    return rs
}

func (d *GAEDatastore) GetKeysOnly(q *datastore.Query) ([]DO) {
    rs := make([]DO, 0)
    keys, err := q.GetAll(d.c.ctx, nil)
    if err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    for i := 0; i < len(keys); i++ {
        rs = append(rs, DO{})
        rs[i]["$key"] = keys[i].Encode()
        rs[i]["$id"] = keys[i].IntID()
        if rs[i]["$id"].(int64) == 0 {
            rs[i]["$id"] = keys[i].StringID()
        }
    }
    return rs
}

func (d *GAEDatastore) GetO(name string, id string, rs interface{}) {
    key := new(datastore.Key)
    if intID, err := strconv.Atoi(id); err == nil {
        key = datastore.NewKey(d.c.ctx, name, "", int64(intID), nil)
    } else {
        key = datastore.NewKey(d.c.ctx, name, id, 0, nil)
    }
    if err := datastore.Get(d.c.ctx, key, &rs); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
}

func (d *GAEDatastore) GetAllO(q *datastore.Query, rs interface{}) {
    _, err := q.GetAll(d.c.ctx, &rs)
    if err != nil {
        d.c.Log.Error("DATASTORE GETALLO %v", err)
    }
}

func (d *GAEDatastore) Put(kind string, data DO) {
    k := datastore.NewIncompleteKey(d.c.ctx, kind, nil)
    key, err := datastore.Put(d.c.ctx, k, &data)
    d.c.Log.Info("%v %v", key, err)
    data["$key"] = key
    // data["_id"] = strconv.FormatInt(key.IntID(), 10)
}

type PageO struct {
    Name string
    Content string
    Status bool
    // LastEdited time.Time
}

func (d *GAEDatastore) PutA() {
    page := PageO{}
    page.Name = "abc"
    page.Content = ""
    page.Status = true

    k := datastore.NewIncompleteKey(d.c.ctx, "PageO", nil)
    _, err := datastore.Put(d.c.ctx, k, &page)
    if err != nil {
        d.c.Log.Error("DATASTORE PUTA %v", err)
    }
}

func (d *GAEDatastore) GetA() []PageO {
    rs := make([]PageO, 0)

    q := datastore.NewQuery("PageO").Project("Name")
    _, err := q.GetAll(d.c.ctx, &rs)
    if err != nil {
        d.c.Log.Error("DATASTORE GETA %v", err)
    }
    // d.c.Log.Error("DATASTORE GETA %v", rs)
    return rs
}

func (d *GAEDatastore) PutO(kind string, data interface{}) {
    k := datastore.NewIncompleteKey(d.c.ctx, kind, nil)
    _, err := datastore.Put(d.c.ctx, k, &data)
    if err != nil {
        d.c.Log.Error("DATASTORE PUTO %v", err)
    }
    //
    // d.c.Log.Info("%v %v", key, err)
    // data["$key"] = key
    // data["_id"] = strconv.FormatInt(key.IntID(), 10)
}

func (d *GAEDatastore) PutOWithKey(key *datastore.Key, data interface{}) {
    if _, err := datastore.Put(d.c.ctx, key, data); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
}

func (d *GAEDatastore) PutWithId(kind string, id string, data DO) {
    key := datastore.NewKey(d.c.ctx, kind, id, 0, nil)
    if _, err := datastore.Put(d.c.ctx, key, &data); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    data["$key"] = key
}

func (d *GAEDatastore) PutWithIntId(kind string, id int, data DO) {
    key := datastore.NewKey(d.c.ctx, kind, "", int64(id), nil)
    if _, err := datastore.Put(d.c.ctx, key, &data); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    data["$key"] = key
}

func (d *GAEDatastore) PutWithKey(key *datastore.Key, data DO) *datastore.Key {
    if _, err := datastore.Put(d.c.ctx, key, &data); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
    data["$key"] = key
    return key
}



func (d *GAEDatastore) Delete(name string, id string) {
    key := datastore.NewKey(d.c.ctx, name, id, 0, nil)
    if err := datastore.Delete(d.c.ctx, key); err != nil {
        d.c.Log.Error("DATASTORE %v", err)
    }
}

func (d *GAEDatastore) DecodeKey(encoded string) (*datastore.Key, error) {
    return datastore.DecodeKey(encoded)
}
*/
