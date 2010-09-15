package kernel

type storeDict map[string]interface{}

type DataStore struct {
	Component
	store storeDict
}

func NewDataStore(name string) *DataStore {
	store := &DataStore{}
	store.Component.comp_name = name
	store.Component.comp_type = "kernel.DataStore"
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

// check implementations match interfaces
var _ = IComponent(&DataStore{})
var _ = IDataStore(&DataStore{})

/* EOF */
