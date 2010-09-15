package kernel

type appMgr struct {
	name string
	jobo string
	
	evtproc IEvtProcessor
	evtsel  IEvtSelector

	mgrs    map[string]IComponentMgr
	svcmgr  svcMgr
	algmgr  algMgr
}

type svcMgr struct {
	services map[string]IService
}

func (mgr *svcMgr) GetComp(n string) IComponent {
	c,ok := mgr.services[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (mgr *svcMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(mgr.services))
	i := 0
	for _,v := range mgr.services {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (mgr *svcMgr) HasComp(n string) bool {
	_,ok := mgr.services[n]
	if !ok {
		mgr.services[n] = nil, false
	}
	return ok
}

func (mgr *svcMgr) AddService(svc string) StatusCode {
	panic("AddService(svc string) not implemented")
	return StatusCode(1)
}

func (mgr *svcMgr) RemoveService(svc string) StatusCode {
	if mgr.HasService(svc) == StatusCode(0) {
		mgr.services[svc] = nil, false
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (mgr *svcMgr) HasService(svc string) StatusCode {
	if mgr.HasComp(svc) {
		return StatusCode(0)
	}
	return StatusCode(1)
}

type algMgr struct {
	algs map[string]IAlgorithm
}

func (mgr *algMgr) GetComp(n string) IComponent {
	c,ok := mgr.algs[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (mgr *algMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(mgr.algs))
	i := 0
	for _,v := range mgr.algs {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (mgr *algMgr) HasComp(n string) bool {
	_,ok := mgr.algs[n]
	if !ok {
		mgr.algs[n] = nil, false
	}
	return ok
}

func (mgr *algMgr) AddAlgorithm(alg IAlgorithm) StatusCode {
	mgr.algs[alg.CompName()] = alg
	return StatusCode(0)
}

func (mgr *algMgr) RemoveAlgorithm(alg IAlgorithm) StatusCode {
	n := alg.CompName()
	if !mgr.HasComp(n) {
		return StatusCode(1)
	}
	mgr.algs[n] = nil, false
	return StatusCode(0)
}

func (mgr *algMgr) HasAlgorithm(algname string) bool {
	return mgr.HasComp(algname)
}

func (app *appMgr) CompType() string {
	return "gaudi.kernel.appMgr"
}

func (app *appMgr) CompName() string {
	return app.name
}

func (app *appMgr) Configure() StatusCode {
	app.evtproc = NewEvtProcessor("evt-proc")
	//app.evtsel  = 

	return StatusCode(0)
}

func (app *appMgr) Initialize() StatusCode {
	println(app.name, "initialize...")
	return StatusCode(0)
}

func (app *appMgr) Start() StatusCode {
	println(app.name, "start...")
	return StatusCode(0)
}

func (app *appMgr) Run() StatusCode {
	println(app.name, "run...")
	sc := app.evtproc.ExecuteRun(10)
	return sc
}

func (app *appMgr) Stop() StatusCode {
	println(app.name, "stop...")
	return StatusCode(0)
}

func (app *appMgr) Finalize() StatusCode {
	println(app.name, "finalize...")
	return StatusCode(0)
}

func (app *appMgr) Terminate() StatusCode {
	println(app.name, "terminate...")
	return StatusCode(0)
}

func (mgr *appMgr) AddAlgorithm(alg IAlgorithm) StatusCode {
	return mgr.algmgr.AddAlgorithm(alg)
}

func (mgr *appMgr) RemoveAlgorithm(alg IAlgorithm) StatusCode {
	return mgr.algmgr.RemoveAlgorithm(alg)
}

func (mgr *appMgr) HasAlgorithm(algname string) bool {
	return mgr.algmgr.HasAlgorithm(algname)
}

func NewAppMgr() IAppMgr {
	appmgr := &appMgr{}
	appmgr.name = "app-mgr"
	appmgr.jobo = "foo.py"

	appmgr.svcmgr = svcMgr{}
	appmgr.svcmgr.services = make(map[string]IService)

	appmgr.algmgr = algMgr{}
	appmgr.algmgr.algs = make(map[string]IAlgorithm)


	appmgr.mgrs = make(map[string]IComponentMgr)
	appmgr.mgrs["svcmgr"] = &appmgr.svcmgr
	appmgr.mgrs["algmgr"] = &appmgr.algmgr

	return appmgr
}

// check implementations match interfaces
var _ = IAlgMgr(&algMgr{})
var _ = IComponentMgr(&algMgr{})

var _ = IComponentMgr(&svcMgr{})
var _ = ISvcMgr(&svcMgr{})

var _ = IComponent(&appMgr{})
var _ = IAlgMgr(&appMgr{})
//var _ = ISvcMgr(&appMgr{})
var _ = IAppMgr(&appMgr{})

/* EOF */
