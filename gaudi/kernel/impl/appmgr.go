package kernel

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
	isvc,ok := components[svc].(IService)
	if !ok {
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
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (mgr *svcMgr) GetService(svc string) IService {
	if mgr.HasService(svc) != StatusCode(0) {
		return mgr.services[svc]
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
	return components[n]
}

func (mgr *appMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(components))
	i := 0
	for _,v := range components {
		comps[i] = v
		i++
	}
	return comps
}

func (mgr *appMgr) HasComp(n string) bool {
	_,ok := components[n]
	if !ok {
		components[n] = nil, false
	}
	return ok
}

func (app *appMgr) Configure() StatusCode {
	app.evtproc = NewEvtProcessor("evt-proc")
	//app.evtsel  = 

	iprop, ok := app.evtproc.(IProperty)
	if ok {
		for _,v := range app.GetProperties() {
			println(app.CompName(), 
				"sending prop[",v.Name,"]=[",v.Value,
				"] to evt-processor...")
			sc := iprop.SetProperty(v.Name, v.Value)
			println(app.CompName(), 
				"sending prop[",v.Name,"]=[",v.Value,
				"] to evt-processor... [",sc,"]")
		}
	}
	return StatusCode(0)
}

func (app *appMgr) Initialize() StatusCode {
	allgood := true
	println(app.name, "initialize...")

	svcs_prop, ok := app.GetProperty("Svcs").([]string)
	if ok {
		for _,svc_name := range svcs_prop {
			isvc := app.GetService(svc_name)
			sc := isvc.InitializeSvc()
			if sc != StatusCode(0) {
				println("** pb initializing [",isvc.CompName(),"] !")
				allgood = false
			}
		}
	}

	algs_prop, ok := app.GetProperty("Algs").([]string)
	if ok {
		for _,alg_name := range algs_prop {
			ialg,isalg := app.GetComp(alg_name).(IAlgorithm)
			if isalg {
				sc := ialg.Initialize()
				if sc != StatusCode(0) {
					println("** pb initializing [",ialg.CompName(),"] !")
					allgood = false
				}
			}
		}
	}
	if allgood {
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (app *appMgr) Start() StatusCode {
	println(app.name, "start...")
	return StatusCode(0)
}

func (app *appMgr) Run() StatusCode {
	println(app.name, "run...")
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
	appmgr.properties.props = make(map[string]interface{})
	appmgr.name = "app-mgr"
	appmgr.jobo = "foo.py"

	appmgr.svcMgr.services = make(map[string]IService)
	appmgr.algMgr.algs = make(map[string]IAlgorithm)


	appmgr.mgrs = make(map[string]IComponentMgr)
	appmgr.mgrs["svcmgr"] = &appmgr.svcMgr
	appmgr.mgrs["algmgr"] = &appmgr.algMgr
	
	components[appmgr.name] = appmgr
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
