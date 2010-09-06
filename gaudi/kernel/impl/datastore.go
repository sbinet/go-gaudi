package kernel

type storeDict map[string]interface{}

type DataStore struct {
	name  string
	store storeDict
}

func NewDataStore(name string) *DataStore {
	store := &DataStore{name:name}
	store.store = make(storeDict)
	return store
}

func (d *DataStore) Put(key string, value *interface{}) bool {
	ok := true
	d.store[key] = value
	return ok
}

func (d *DataStore) Has(key string) bool {
	_,ok := d.store[key]
	return ok
}

func (d *DataStore) Get(key string) (chan *interface{}, bool) {
	out := make(chan *interface{})
	v,ok := d.store[key]
	if ok {
		out <- &v
	} else {
		out <- nil
	}
	return out, ok
}

