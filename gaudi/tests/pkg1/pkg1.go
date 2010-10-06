// test package 'pkg1'
package pkg1

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
	self.MsgInfo("== execute == [ctx:%v]\n", ctx.Idx())
	return kernel.StatusCode(0)
}

func (self *alg1) Finalize() kernel.StatusCode {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- alg2 ---

type alg2 struct {
	kernel.Algorithm
}

func (self *alg2) Initialize() kernel.StatusCode {
	self.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (self *alg2) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgInfo("~~ execute ~~ [ctx:%v]\n", ctx.Idx())
	return kernel.StatusCode(0)
}

func (self *alg2) Finalize() kernel.StatusCode {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- svc1 ---
type svc1 struct {
	kernel.Service
}

func (self *svc1) InitializeSvc() kernel.StatusCode {
	self.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (self *svc1) FinalizeSvc() kernel.StatusCode {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- tool1 ---
type tool1 struct {
	kernel.AlgTool
}

func (self *tool1) InitializeTool() kernel.StatusCode {
	self.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (self *tool1) FinalizeTool() kernel.StatusCode {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// check matching interfaces
var _ = kernel.IComponent(&alg1{})
var _ = kernel.IAlgorithm(&alg1{})
//var _ = kernel.Algorithm(&alg1{})

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "alg1":
		self := &alg1{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	case "alg2":
		self := &alg2{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	case "svc1":
		self := &svc1{}
		_ = kernel.NewSvc(&self.Service,t,n)
		kernel.RegisterComp(self)
		return self
	case "tool1":
		self := &tool1{}
		_ = kernel.NewTool(&self.AlgTool,t,n, nil)
		kernel.RegisterComp(self)
		return self
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
