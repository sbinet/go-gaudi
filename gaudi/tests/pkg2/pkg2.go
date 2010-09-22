// test package 'pkg2'
package pkg2

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
	a.MsgInfo("== execute == [%v]\n", ctx)
	return kernel.StatusCode(0)
}

func (a *alg1) Finalize() kernel.StatusCode {
	a.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- svc2 ---
type svc2 struct {
	kernel.Service
}

func (s *svc2) InitializeSvc() kernel.StatusCode {
	s.MsgInfo("~~ initialize ~~\n")
	return kernel.StatusCode(0)
}

func (s *svc2) FinalizeSvc() kernel.StatusCode {
	s.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "alg1":
		c := &alg1{}
		_ = kernel.NewAlg(&c.Algorithm,t,n)
		kernel.RegisterComp(c)
		return c
	case "svc2":
		c := &svc2{}
		_ = kernel.NewSvc(&c.Service,t,n)
		kernel.RegisterComp(c)
		kernel.RegisterComp(c)
		return c
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
