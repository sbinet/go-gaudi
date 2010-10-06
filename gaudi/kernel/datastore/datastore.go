// reference implementation of an IDataStore
package datastore

import "gaudi/kernel"

// --- datastore helper ---

type datastore struct {
	store *kernel.DataStore //map[string]interface{}
}

func (self *datastore) Put(key string, value interface{}) {
	(*self.store)[key] = value
}

func (self *datastore) Get(key string) interface{} {
	value, ok := (*self.store)[key]
	if ok {
		return value
	}
	return nil
}

func (self *datastore) Has(key string) bool {
	_, ok := (*self.store)[key]
	if !ok {
		(*self.store)[key] = nil, false
	}
	return ok
}

func (self *datastore) ClearStore() kernel.StatusCode {
	for k,_ := range (*self.store) {
		(*self.store)[k] = nil, false
	}
	return kernel.StatusCode(0)
}

// --- datastore service ---

type datastoresvc struct {
	kernel.Service
}

func (self *datastoresvc) InitializeSvc() kernel.StatusCode {
	self.MsgInfo("~~ initialize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) FinalizeSvc() kernel.StatusCode {
	self.MsgInfo("~~ finalize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) Store(ctx kernel.IEvtCtx) kernel.IDataStore {
	return &datastore{ctx.Store()}
}

// check matching interfaces
var _ = kernel.IDataStore(&datastore{})
var _ = kernel.IComponent(&datastoresvc{})
var _ = kernel.IService(&datastoresvc{})
var _ = kernel.IProperty(&datastoresvc{})

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "datastoresvc":
		self := &datastoresvc{}
		//self.stores = make([]datastore, 1)
		_ = kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)
		return self
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}

/* EOF */