// test package 'pkg2'
package pkg2

import "gaudi/kernel"

// --- alg1 ---

type alg1 struct {
	kernel.Algorithm
}

func (self *alg1) Initialize() kernel.StatusCode {
	self.MsgInfo("== initialize ==\n")
	return kernel.StatusCode(0)
}

func (self *alg1) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgInfo("== execute == [%v]\n", ctx)
	return kernel.StatusCode(0)
}

func (self *alg1) Finalize() kernel.StatusCode {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- svc2 ---
type svc2 struct {
	kernel.Service
}

func (self *svc2) InitializeSvc() kernel.StatusCode {
	self.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (self *svc2) FinalizeSvc() kernel.StatusCode {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- alg adder ---
type alg_adder struct {
	kernel.Algorithm
	evtstore kernel.IDataStore
}

func (self *alg_adder) Initialize() kernel.StatusCode {
	self.MsgInfo("== initialize ==\n")
	self.MsgInfo("retrieving evt-store...\n")
	svcloc := kernel.GetSvcLocator()
	self.evtstore = svcloc.GetService("evt-store").(kernel.IDataStore)
	if self.evtstore == nil {
		self.MsgError("could not retrieve evt-store !\n")
		return kernel.StatusCode(1)
	}
	self.MsgInfo("retrieving evt-store... [ok]\n")
	return kernel.StatusCode(0)
}

func (self *alg_adder) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgInfo("== execute == [%v]\n", ctx)
	//self.evtstore.Put("njets", 2)
	return kernel.StatusCode(0)
}

func (self *alg_adder) Finalize() kernel.StatusCode {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "alg1":
		self := &alg1{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	case "svc2":
		self := &svc2{}
		_ = kernel.NewSvc(&self.Service,t,n)
		kernel.RegisterComp(self)
		return self
	case "alg_adder":
		self := &alg_adder{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
