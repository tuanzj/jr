package jr

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type XMGO struct {
	ctx *Context
	db  *mgo.Database
}

// type BSONObjectId bson.ObjectId

// type MGORoot struct {
// 	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
// }

// type Brand struct {
// 	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
// 	Name      string        `json:"name"`
// 	EventId   int           `json:"event_id"`
// 	GiftType  int           `json:"gift_type"`
// 	Percent   int           `json:"percent"`
// 	Status    int           `json:"status"`
// 	Owner     int           `json:"owner"`
// 	Gift      string        `json:"gift"`
// 	ProjectId string        `json:"project_id"`
// }

// type Raw struct {
// 	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
// 	Raw     string        `json:"raw"`
// 	Content string        `json:"content"`
// }

func (d *XMGO) Connect(connect ...string) error {
	session, err := mgo.Dial(connect[0] + "://" + connect[1])
	if err != nil {
		return err
	}
	d.db = session.DB(connect[2])
	return nil
}

func (d *XMGO) Disconnect() error {
	d.db.Session.Close()
	return nil
}

func (d *XMGO) Find(collection string, id interface{}, dst interface{}) error {
	query := d.getQuery(collection, bson.M{"_id": id})
	err := query.One(dst)
	if err != nil {
		return err
	}
	return nil
}
func (d *XMGO) getQuery(collection string, selector map[string]interface{}) *mgo.Query {
	doc := d.db.C(collection)
	query := make(map[string]interface{}, 0)
	keyword := []string{
		"$limit", "$skip", "$sort", "$project",
	}
	for k, v := range selector {
		exists, _ := InArray(k, keyword)
		if !exists {
			query[k] = v
		}
	}
	q := doc.Find(query)
	for k, v := range selector {
		switch k {
		case "$limit":
			if limit, ok := v.(int); ok {
				q = q.Limit(limit)
			}
		case "$sort":
			if order, ok := v.(string); ok {
				q = q.Sort(order)
			} else if orders, ok := v.([]string); ok {
				q = q.Sort(orders...)
			}
		case "$skip":
			if skip, ok := v.(int); ok {
				q = q.Skip(skip)
			}
		case "$project":
			if projects, ok := v.([]string); ok {
				selc := make(bson.M, 0)
				for i := 0; i < len(projects); i++ {
					selc[projects[i]] = 1
				}
				q = q.Select(selc)
			}
		}
	}
	return q
}

func (d *XMGO) FindAll(collection string, selector map[string]interface{}, dst interface{}) error {
	query := d.getQuery(collection, selector)
	err := query.All(dst)
	if err != nil {
		return err
	}
	// num, err := tb.Find(bson.M{
	// 	"mac": mac,
	// 	"time": bson.M{
	// 		"$gte": time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, vn),
	// 		"$lt":  time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, vn),
	// 	}}).
	return nil
}

func (d *XMGO) InsertAll(collection string, docs ...interface{}) error {
	return nil
}

func (d *XMGO) Insert(collection string, doc interface{}, id ...interface{}) error {
	coll := d.db.C(collection)
	err := coll.Insert(doc)
	if err != nil {
		return err
	}
	return nil
}

func (d *XMGO) Update(collection string, doc interface{}, id interface{}) error {
	// var idd interface{}
	// if len(id) > 0 {
	// 	idd = id[0]
	// } else {
	// 	// idd = doc.Id
	// }
	coll := d.db.C(collection)
	err := coll.Update(bson.M{"_id": id}, doc) //bson.M{"$set": doc}
	if err != nil {
		return err
	}
	return nil
}

func (d *XMGO) UpdateAll(collection string, selector map[string]interface{}, doc interface{}) error {
	return nil
}

func (d *XMGO) Upsert(collection string, doc interface{}) error {
	// if id, ok := doc.Id; ok {

	// }
	// id := doc.(IGAEDO).GetID()
	// var err error
	// if id == "" {
	// 	err = d.Insert(collection, doc)
	// } else {
	// 	err = d.Update(collection, doc, id)
	// }
	// if err != nil {
	// 	return err
	// }
	// return nil

	// coll := d.db.C(collection)
	// err := coll.Upsert(bson.M{"_id": id}, bson.M{"$set": doc})
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (d *XMGO) Remove(collection string, id interface{}) error {
	// var idd interface{}
	// if len(id) > 0 {
	// 	idd = id[0]
	// } else {
	// 	// idd = doc.Id
	// }
	coll := d.db.C(collection)
	err := coll.Remove(bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (d *XMGO) RemoveAll(collection string, selector map[string]interface{}) error {
	return nil
}
