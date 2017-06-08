package jr

import (
	"google.golang.org/appengine/datastore"
)

type DO map[string]interface{}

/*
func (d *DO) Load(ch <-chan datastore.Property) error {
    // Note: you might want to clear current values from the map or create a new map
    for p := range ch { // Read until channel is closed
        (*d)[p.Name] = p.Value
    }
    return nil
}

func (d *DO) Save(ch chan<- datastore.Property) error {
    defer close(ch) // Channel must be closed
    for k, v := range *d {
        ch <- datastore.Property{Name: k, Value: v}
    }
    return nil
}
*/
func (d *DO) Load(p []datastore.Property) error {
	// CX.Log.Print("CX Load", p)
	for i := range p {
		// CX.Log.Print("CX Load a", p[i].Value)
		(*d)[p[i].Name] = p[i].Value
	}
	return nil
	/*
	   x := *d
	   for _, p := range ch {
	       if reflect.TypeOf(p.Value) == reflect.TypeOf(datastore.indexValue{}) && reflect.ValueOf(p.Value).Kind() == reflect.Struct {
	           y := p.Value.(datastore.IndexValue).Value
	           switch {
	           case y.Int64Value != nil:
	               x[p.Name] = y.GetInt64Value()
	           case y.Pointvalue != nil:
	               point := y.GetPointvalue()
	               x[p.Name] = appengine.GeoPoint{Lat: point.GetX(), Lng: point.GetY()}
	           case y.StringValue != nil:
	               x[p.Name] = y.GetStringValue()
	           case y.BooleanValue != nil:
	               x[p.Name] = y.GetBooleanValue()
	           case y.DoubleValue != nil:
	               x[p.Name] = y.GetDoubleValue()
	           }
	       } else {
	           x[p.Name] = p.Value
	       }
	   }
	   return nil
	*/
}

func (d *DO) Save() ([]datastore.Property, error) {
	a := make([]datastore.Property, 0)
	// CX.Log.Print("CX Save", d)
	for k, v := range *d {
		a = append(a, datastore.Property{Name: k, Value: v})
	}
	return a, nil
}

func (d *DO) SetID(id string) {
	(*d)["ID"] = id
}

func (d *DO) GetID() string {
	return (*d)["ID"].(string)
}

/*
func (l *PropertyList) Load(p []Property) error {
    *l = append(*l, p...)
    return nil
}

// Save saves all of l's properties as a slice or Properties.
func (l *PropertyList) Save() ([]Property, error) {
    return *l, nil
}
*/
