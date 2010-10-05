// reference implementation of an IDataStore
package datastore

import "gaudi/kernel"

type datastore map[string]interface{}

// --- datastore service ---

type datastoresvc struct {
	kernel.Service
	store datastore
}

func (self *datastoresvc) InitializeSvc() kernel.StatusCode {
	self.MsgInfo("~~ initialize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) FinalizeSvc() kernel.StatusCode {
	self.MsgInfo("~~ finalize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

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

/*

type IDataStoreMgr interface {
	Store(ctx IEvtCtx) chan IDataStore
}
*/

// check matching interfaces
var _ = kernel.IComponent(&datastoresvc{})
var _ = kernel.IService(&datastoresvc{})
var _ = kernel.IProperty(&datastoresvc{})
var _ = kernel.IDataStore(&datastoresvc{})

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "datastoresvc":
		self := &datastoresvc{}
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
