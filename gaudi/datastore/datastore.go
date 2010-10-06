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
	//self.store = make(kernel.DataStore)
	return kernel.StatusCode(0)
}

// --- datastore service ---

type datastoresvc struct {
	kernel.Service
	stores []datastore
}

func (self *datastoresvc) InitializeSvc() kernel.StatusCode {
	self.MsgInfo("~~ initialize [datastore svc] nstores: %v ~~\n", 
		len(self.stores))
	return kernel.StatusCode(0)
}

func (self *datastoresvc) FinalizeSvc() kernel.StatusCode {
	self.MsgInfo("~~ finalize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) Store(ctx kernel.IEvtCtx) kernel.IDataStore {
	/*
	nstores := len(self.stores)
	idx := ctx.Idx() % nstores
	if idx < nstores {
		self.MsgInfo("==> ctx=%03v idx=%03v nstores=%v\n", ctx.Idx(), idx, nstores)
		return &self.stores[idx]
	}
	return nil
	 */
	return &datastore{ctx.Store()}
}

func (self *datastoresvc) SetNbrStreams(n int) kernel.StatusCode {
	//self.nstores = n
	self.stores = make([]datastore, n)
	return kernel.StatusCode(0)
}

/*
func (self *datastoresvc) Put(key string, value interface{}) {
	self.store[key] = value
}

func (self *datastoresvc) Get(key string) interface{} {
	value, ok := self.store[key]
	if ok {
		return value
	}
	return nil
}

func (self *datastoresvc) Has(key string) bool {
	_, ok := self.store[key]
	if !ok {
		self.store[key] = nil, false
	}
	return ok
}
*/

/*
func (self *datastoresvc) ClearStore() kernel.StatusCode {
	self.store = make(datastore)
	return kernel.StatusCode(0)
}
*/

/*

type IDataStoreMgr interface {
	Store(ctx IEvtCtx) chan IDataStore
}
*/

// check matching interfaces
var _ = kernel.IDataStore(&datastore{})
var _ = kernel.IDataStoreClearer(&datastore{})
var _ = kernel.IComponent(&datastoresvc{})
var _ = kernel.IService(&datastoresvc{})
var _ = kernel.IProperty(&datastoresvc{})
//var _ = kernel.IDataStore(&datastoresvc{})
//var _ = kernel.IDataStoreClearer(&datastoresvc{})

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "datastoresvc":
		self := &datastoresvc{}
		//self.stores = make([]datastore, 1)
		_ = kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)
		self.SetNbrStreams(1)
		return self
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}

/* EOF */
