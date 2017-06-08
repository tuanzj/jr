package jr

// import (
// 	"strings"
// )

type IDO interface {
	SetId(string)
	Id() string
}

type IDB interface {
	// NewKey(colection string, id string) IKey
	// NewIncompleteKey(colection string) IKey
	// DecodeKey(encoded string) (IKey, error)
	// NewQuery(colection string) IQuery
	// FindByKey(key IKey, dst interface{}) error
	// FindOne(selectore IQuery, dst interface{}) error
	// Find(selectore IQuery, dst interface{}) error
	// Insert(docs ...interface{}) error
	// Update(selector IQuery, doc interface{}) error
	// UpdateAll(selector IQuery, doc interface{}) error
	// Upsert(selector IQuery, doc interface{}) error
	// Remove(selector IQuery) error
	// RemoveAll(selector IQuery) error
	Find(collection string, id interface{}, dst interface{}) error
	FindAll(collection string, selector map[string]interface{}, dst interface{}) error
	InsertAll(collection string, docs ...interface{}) error
	Insert(collection string, doc interface{}, id ...interface{}) error
	Update(collection string, doc interface{}, id interface{}) error
	UpdateAll(collection string, selector map[string]interface{}, doc interface{}) error
	Upsert(collection string, doc interface{}) error
	Remove(collection string, id interface{}) error
	RemoveAll(collection string, selector map[string]interface{}) error
	Disconnect() error
	Connect(connect ...string) error
}

/*
type IQuery interface {
    Filter(patern string, value interface{}) IQuery
    Order(patern string) IQuery
    Limit(patern string) IQuery
    Project(fieldNames ...string) IQuery
    // Distinct() IQuery
    KeysOnly() IQuery
    Offset() IQuery
    Start(c ICursor) IQuery
    End(c ICursor) IQuery
    // Count() int
}

type ICursor interface {

}

type IKey interface {

}
*/
func newDB(ctx *Context, connect ...string) IDB {
	if len(connect) > 1 {
		switch connect[0] {
		case "mongodb":
			o := &XMGO{
				ctx: ctx,
			}
			o.Connect(connect...)
			return o
		case "mysql":
		}
	} else if ctx.isGAE() {
		return &GAEDatastore{
			c: ctx,
		}
	}
	return nil
}
