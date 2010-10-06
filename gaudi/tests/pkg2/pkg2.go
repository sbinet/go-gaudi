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
	self.MsgInfo("== execute == [%v]\n", ctx.Idx())
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

type simple_counter struct {
	cnt int
}

// --- alg adder ---
type alg_adder struct {
	kernel.Algorithm
	val float
	cnt_key string
}

func (self *alg_adder) Initialize() kernel.StatusCode {
	self.MsgInfo("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	
	/*
	self.MsgInfo("retrieving evt-store...\n")
	svcloc := kernel.GetSvcLocator()
	self.evtstore = svcloc.GetService("evt-store").(kernel.IDataStore)
	if self.evtstore == nil {
		self.MsgError("could not retrieve evt-store !\n")
		return kernel.StatusCode(1)
	}
	self.MsgInfo("retrieving evt-store... [ok]\n")
	 */
	self.MsgInfo("--> val: %v\n", self.val)
	self.val = self.GetProperty("Val").(float)
	self.MsgInfo("--> val: %v\n", self.val)
	self.cnt_key = self.GetProperty("SimpleCounter").(string)
	self.MsgInfo("--> cnt: %v\n", self.cnt_key)
	return kernel.StatusCode(0)
}

func (self *alg_adder) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgInfo("== execute == [%v]\n", ctx.Idx())

	njets := 1 + ctx.Idx()
	val := self.val + 1
	store := self.EvtStore(ctx)

	if store.Has("njets") {
		njets += store.Get("njets").(int)
	}
	store.Put("njets", njets)

	if store.Has("ptjets") {
		val += store.Get("ptjets").(float)
	}
	store.Put("ptjets", val)

	cnt := &simple_counter{0}
	if store.Has(self.cnt_key) {
		cnt = store.Get(self.cnt_key).(*simple_counter)
	}
	cnt.cnt += 1
	store.Put(self.cnt_key, cnt)
	return kernel.StatusCode(0)
}

func (self *alg_adder) Finalize() kernel.StatusCode {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- alg_dumper ---
type alg_dumper struct {
	kernel.Algorithm
	njets_key string
	ptjets_key string
	cnt_key string
	cnt_val int
}

func (self *alg_dumper) Initialize() kernel.StatusCode {
	self.MsgInfo("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	self.njets_key  = self.GetProperty("NbrJets").(string)
	self.ptjets_key = self.GetProperty("PtJets").(string)
	self.cnt_key = self.GetProperty("SimpleCounter").(string)
	self.cnt_val = self.GetProperty("ExpectedValue").(int)
	return kernel.StatusCode(0)
}

func (self *alg_dumper) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgInfo("== execute == [%v]\n", ctx.Idx())

	store := self.EvtStore(ctx)
	njets  := store.Get(self.njets_key).(int)
	ptjets := store.Get(self.ptjets_key).(float)
	cnt    := store.Get(self.cnt_key).(*simple_counter)
	allgood := "ERR"
	if self.cnt_val == cnt.cnt {
		allgood = "OK"
	}
	self.MsgInfo("[ctx:%03v] njets: %03v ptjets: %8.3v [%s]\n", 
		ctx.Idx(), njets, ptjets, allgood)

	return kernel.StatusCode(0)
}

func (self *alg_dumper) Finalize() kernel.StatusCode {
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

		self.DeclareProperty("Val", -99.)
		self.DeclareProperty("SimpleCounter", "cnt")
		return self

	case "alg_dumper":
		self := &alg_dumper{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)

		self.DeclareProperty("SimpleCounter", "cnt")
		self.DeclareProperty("ExpectedValue", -1)
		self.DeclareProperty("NbrJets", "njets")
		self.DeclareProperty("PtJets",  "ptjets")
		return self

	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
