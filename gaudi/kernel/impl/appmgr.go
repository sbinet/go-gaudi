package kernel

//import "fmt"

///////////////////////////////////////////////////////////////////////////////
// svc-mgr

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
	isvc,ok := g_compsdb[svc].(IService)
	if !ok {
		//fmt.Printf("** AddService(%s) FAILED !\n", svc)
		return StatusCode(1)
	}
	mgr.services[svc] = isvc
	return StatusCode(0)
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
		//fmt.Printf(":: HasService(%s) - true\n", svc)
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (mgr *svcMgr) GetService(svc string) IService {
	if mgr.HasService(svc).IsSuccess() {
		//fmt.Printf("-- GetService(%s)...\n", svc)
		isvc := mgr.services[svc]
		//fmt.Printf("-- GetService(%s)... [done]\n", svc)
		return isvc
	}
	return nil
}

func (mgr *svcMgr) GetServices() []IService {
	svcs := make([]IService, len(mgr.services))
	i := 0
	for _,v := range mgr.services {
		svcs[i] = v
		i++
	}
	return svcs
}

func (mgr *svcMgr) ExistsService(svc string) bool {
	return mgr.HasService(svc) == StatusCode(0)
}


//////////////////////////////////////////////////////////////////////////////
// alg-mgr

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

//////////////////////////////////////////////////////////////////////////////
// app-mgr

type appMgr struct {
	properties
	msgstream
	svcMgr
	algMgr

	name string
	jobo string
	
	evtproc IEvtProcessor
	evtsel  IEvtSelector

	mgrs    map[string]IComponentMgr
}

func (app *appMgr) CompType() string {
	return "gaudi.kernel.appMgr"
}

func (app *appMgr) CompName() string {
	return app.name
}

func (mgr *appMgr) GetComp(n string) IComponent {
	if !mgr.HasComp(n) {
		return nil
	}
	return g_compsdb[n]
}

func (mgr *appMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(g_compsdb))
	i := 0
	for _,v := range g_compsdb {
		comps[i] = v
		i++
	}
	return comps
}

func (mgr *appMgr) HasComp(n string) bool {
	_,ok := g_compsdb[n]
	if !ok {
		g_compsdb[n] = nil, false
	}
	return ok
}

func (app *appMgr) Configure() StatusCode {
	app.evtproc = NewEvtProcessor("evt-proc")
	//app.evtsel  = 

	return StatusCode(0)
}

func (app *appMgr) Initialize() StatusCode {
	allgood := true
	app.MsgInfo("initialize...\n")

	app.MsgVerbose("components-map: %v\n", g_compsdb)
	svcs_prop, ok := app.GetProperty("Svcs").([]string)
	if ok {
		app.MsgInfo("svcs...\n")
		for _,svc_name := range svcs_prop {
			isvc := app.GetService(svc_name)
			sc := isvc.InitializeSvc()
			if sc != StatusCode(0) {
				app.MsgError("pb initializing [%s] !\n",isvc.CompName())
				allgood = false
			}
		}
		_ = app.evtproc.(IService).InitializeSvc()
		app.MsgInfo("svcs... [done]\n")
	}

	algs_prop, ok := app.GetProperty("Algs").([]string)
	if ok {
		app.MsgInfo("algs...\n")
		for _,alg_name := range algs_prop {
			ialg,isalg := app.GetComp(alg_name).(IAlgorithm)
			if isalg {
				sc := ialg.Initialize()
				if sc != StatusCode(0) {
					app.MsgError("pb initializing [%s] !\n",ialg.CompName())
					allgood = false
				} else {
					app.MsgDebug("correctly initialized [%T/%s]\n",
						ialg, ialg.CompName())
				}
			}
		}
		app.MsgInfo("algs... [done]\n")
	}
	if allgood {
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (app *appMgr) Start() StatusCode {
	app.MsgInfo("start...\n")
	return StatusCode(0)
}

func (app *appMgr) Run() StatusCode {
	app.MsgInfo("run...\n")
	// init
	sc := app.Initialize()
	if !sc.IsSuccess() {
		return sc
	}

	// start
	sc = app.Start()
	if !sc.IsSuccess() {
		return sc
	}

	// evtloop-run
	sc = app.evtproc.ExecuteRun(10)
	if !sc.IsSuccess() {
		return sc
	}

	// stop
	sc = app.Stop()
	if !sc.IsSuccess() {
		return sc
	}

	// fini
	sc = app.Finalize()
	if !sc.IsSuccess() {
		return sc
	}

	return app.Terminate()
}

func (app *appMgr) Stop() StatusCode {
	app.MsgInfo("stop...\n")
	return StatusCode(0)
}

func (app *appMgr) Finalize() StatusCode {
	app.MsgInfo("finalize...\n")
	allgood := true

	svcs_prop, ok := app.GetProperty("Svcs").([]string)
	if ok {
		app.MsgInfo("svcs...\n")
		for _,svc_name := range svcs_prop {
			isvc := app.GetService(svc_name)
			sc := isvc.FinalizeSvc()
			if sc != StatusCode(0) {
				app.MsgError("pb finalizing [%s] !\n",isvc.CompName())
				allgood = false
			}
		}
		_ = app.evtproc.(IService).FinalizeSvc()
		app.MsgInfo("svcs... [done]\n")
	}

	algs_prop, ok := app.GetProperty("Algs").([]string)
	if ok {
		app.MsgInfo("algs...\n")
		for _,alg_name := range algs_prop {
			ialg,isalg := app.GetComp(alg_name).(IAlgorithm)
			if isalg {
				sc := ialg.Finalize()
				if sc != StatusCode(0) {
					app.MsgError("pb finalizing [%s] !\n",ialg.CompName())
					allgood = false
				} else {
					app.MsgDebug("correctly finalized [%T/%s]\n",
						ialg, ialg.CompName())
				}
			}
		}
		app.MsgInfo("algs... [done]\n")
	}
	if allgood {
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (app *appMgr) Terminate() StatusCode {
	app.MsgInfo("terminate...\n")
	return StatusCode(0)
}

func NewAppMgr() IAppMgr {
	appmgr := &appMgr{}
	appmgr.properties.props = make(map[string]interface{})
	appmgr.name = "app-mgr"
	appmgr.jobo = "foo.py"
	appmgr.msgstream = msgstream{name:appmgr.name, level:LVL_INFO}

	appmgr.svcMgr.services = make(map[string]IService)
	appmgr.algMgr.algs = make(map[string]IAlgorithm)


	appmgr.mgrs = make(map[string]IComponentMgr)
	appmgr.mgrs["svcmgr"] = &appmgr.svcMgr
	appmgr.mgrs["algmgr"] = &appmgr.algMgr
	
	g_compsdb[appmgr.name] = appmgr

	// completing bootstrap
	g_isvcloc = appmgr

	return appmgr
}

// check implementations match interfaces
var _ = IAlgMgr(&algMgr{})
var _ = IComponentMgr(&algMgr{})

var _ = IComponentMgr(&svcMgr{})
var _ = ISvcMgr(&svcMgr{})
var _ = ISvcLocator(&svcMgr{})

var _ = IComponent(&appMgr{})
var _ = IComponentMgr(&appMgr{})
var _ = IAlgMgr(&appMgr{})
var _ = ISvcMgr(&appMgr{})
var _ = ISvcLocator(&appMgr{})
var _ = IAppMgr(&appMgr{})
var _ = IProperty(&appMgr{})


/* EOF */
