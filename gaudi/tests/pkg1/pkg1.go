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

// --- alg2 ---

type Alg2 struct {
	kernel.Algorithm
}

func (a *Alg2) Initialize() kernel.StatusCode {
	println(a.CompName(), "~~ initialize ~~")
	return kernel.StatusCode(0)
}

func (a *Alg2) Execute() kernel.StatusCode {
	println(a.CompName(), "~~ execute ~~")
	return kernel.StatusCode(0)
}

func (a *Alg2) Finalize() kernel.StatusCode {
	println(a.CompName(), "~~ finalize ~~")
	return kernel.StatusCode(0)
}

// --- svc1 ---
type Svc1 struct {
	kernel.Service
}

func (s *Svc1) Initialize() kernel.StatusCode {
	println(s.CompName(), "~~ initialize ~~")
	return kernel.StatusCode(0)
}

func (s *Svc1) Finalize() kernel.StatusCode {
	println(s.CompName(), "~~ finalize ~~")
	return kernel.StatusCode(0)
}

// --- tool1 ---
type Tool1 struct {
	kernel.AlgTool
}

func (t *Tool1) Initialize() kernel.StatusCode {
	println(t.CompName(), "~~ initialize ~~")
	return kernel.StatusCode(0)
}

func (t *Tool1) Finalize() kernel.StatusCode {
	println(t.CompName(), "~~ finalize ~~")
	return kernel.StatusCode(0)
}


/* EOF */
