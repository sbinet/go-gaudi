// test package 'pkg1'
package pkg1

import "gaudi/kernel"

// --- alg1 ---

type Alg1 struct {
	kernel.Algorithm
}

func (a *Alg1) Initialize() kernel.StatusCode {
	println(a.CompName(), "== initialize ==")
	return kernel.StatusCode(0)
}

func (a *Alg1) Execute() kernel.StatusCode {
	println(a.CompName(), "== execute ==")
	return kernel.StatusCode(0)
}

func (a *Alg1) Finalize() kernel.StatusCode {
	println(a.CompName(), "== finalize ==")
	return kernel.StatusCode(0)
}

// --- svc2 ---
type Svc2 struct {
	kernel.Service
}

func (s *Svc2) Initialize() kernel.StatusCode {
	println(s.CompName(), "~~ initialize ~~")
	return kernel.StatusCode(0)
}

func (s *Svc2) Finalize() kernel.StatusCode {
	println(s.CompName(), "~~ finalize ~~")
	return kernel.StatusCode(0)
}

/* EOF */
