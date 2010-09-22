// test package 'pkg1'
package pkg1

import "gaudi/kernel"

// --- alg1 ---

type alg1 struct {
	kernel.Algorithm
}

func (a *alg1) Initialize() kernel.StatusCode {
	a.MsgInfo("== initialize ==\n")
	return kernel.StatusCode(0)
}

func (a *alg1) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	a.MsgInfo("== execute == [ctx:%v]\n", ctx)
	return kernel.StatusCode(0)
}

func (a *alg1) Finalize() kernel.StatusCode {
	a.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- alg2 ---

type alg2 struct {
	kernel.Algorithm
}

func (a *alg2) Initialize() kernel.StatusCode {
	a.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (a *alg2) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	a.MsgInfo("~~ execute ~~ [ctx:%v]\n", ctx)
	return kernel.StatusCode(0)
}

func (a *alg2) Finalize() kernel.StatusCode {
	a.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- svc1 ---
type svc1 struct {
	kernel.Service
}

func (s *svc1) InitializeSvc() kernel.StatusCode {
	s.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (s *svc1) FinalizeSvc() kernel.StatusCode {
	s.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- tool1 ---
type tool1 struct {
	kernel.AlgTool
}

func (t *tool1) InitializeTool() kernel.StatusCode {
	t.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (t *tool1) FinalizeTool() kernel.StatusCode {
	t.MsgInfo("~~ finalize ~~\n")
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
		c := &alg1{}
		_ = kernel.NewAlg(&c.Algorithm,t,n)
		kernel.RegisterComp(c)
		return c
	case "alg2":
		c := &alg2{}
		_ = kernel.NewAlg(&c.Algorithm,t,n)
		kernel.RegisterComp(c)
		return c
	case "svc1":
		c := &svc1{}
		_ = kernel.NewSvc(&c.Service,t,n)
		kernel.RegisterComp(c)
		return c
	case "tool1":
		c := &tool1{}
		_ = kernel.NewTool(&c.AlgTool,t,n, nil)
		kernel.RegisterComp(c)
		return c
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
