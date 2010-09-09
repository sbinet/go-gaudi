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

func (app *appMgr) CompType() string {
	return "gaudi.kernel.appMgr"
}

func (app *appMgr) CompName() string {
	return app.name
}

func (app *appMgr) Configure() StatusCode {
	app.evtproc = NewEvtProcessor("evt-proc")
	//app.evtsel  = 

	return StatusCode(-1)
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

