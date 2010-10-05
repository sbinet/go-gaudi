package kernel

//import "fmt"

///////////////////////////////////////////////////////////////////////////////
// svc-mgr

type svcMgr struct {
	services map[string]IService
}

func (self *svcMgr) GetComp(n string) IComponent {
	c,ok := self.services[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (self *svcMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(self.services))
	i := 0
	for _,v := range self.services {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (self *svcMgr) HasComp(n string) bool {
	_,ok := self.services[n]
	if !ok {
		self.services[n] = nil, false
	}
	return ok
}

func (self *svcMgr) AddService(svc string) StatusCode {
	isvc,ok := g_compsdb[svc].(IService)
	if !ok {
		//fmt.Printf("** AddService(%s) FAILED !\n", svc)
		return StatusCode(1)
	}
	self.services[svc] = isvc
	return StatusCode(0)
}

func (self *svcMgr) RemoveService(svc string) StatusCode {
	if self.HasService(svc) == StatusCode(0) {
		self.services[svc] = nil, false
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *svcMgr) HasService(svc string) StatusCode {
	if self.HasComp(svc) {
		//fmt.Printf(":: HasService(%s) - true\n", svc)
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *svcMgr) GetService(svc string) IService {
	if self.HasService(svc).IsSuccess() {
		//fmt.Printf("-- GetService(%s)...\n", svc)
		isvc := self.services[svc]
		//fmt.Printf("-- GetService(%s)... [done]\n", svc)
		return isvc
	}
	return nil
}

func (self *svcMgr) GetServices() []IService {
	svcs := make([]IService, len(self.services))
	i := 0
	for _,v := range self.services {
		svcs[i] = v
		i++
	}
	return svcs
}

func (self *svcMgr) ExistsService(svc string) bool {
	return self.HasService(svc) == StatusCode(0)
}


//////////////////////////////////////////////////////////////////////////////
// alg-mgr

type algMgr struct {
	algs map[string]IAlgorithm
}

func (self *algMgr) GetComp(n string) IComponent {
	c,ok := self.algs[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (self *algMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(self.algs))
	i := 0
	for _,v := range self.algs {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (self *algMgr) HasComp(n string) bool {
	_,ok := self.algs[n]
	if !ok {
		self.algs[n] = nil, false
	}
	return ok
}

func (self *algMgr) AddAlgorithm(alg IAlgorithm) StatusCode {
	self.algs[alg.CompName()] = alg
	return StatusCode(0)
}

func (self *algMgr) RemoveAlgorithm(alg IAlgorithm) StatusCode {
	n := alg.CompName()
	if !self.HasComp(n) {
		return StatusCode(1)
	}
	self.algs[n] = nil, false
	return StatusCode(0)
}

func (self *algMgr) HasAlgorithm(algname string) bool {
	return self.HasComp(algname)
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

func (self *appMgr) CompType() string {
	return "gaudi.kernel.appMgr"
}

func (self *appMgr) CompName() string {
	return self.name
}

func (self *appMgr) GetComp(n string) IComponent {
	if !self.HasComp(n) {
		return nil
	}
	return g_compsdb[n]
}

func (self *appMgr) GetComps() []IComponent {
	comps := make([]IComponent, len(g_compsdb))
	i := 0
	for _,v := range g_compsdb {
		comps[i] = v
		i++
	}
	return comps
}

func (self *appMgr) HasComp(n string) bool {
	_,ok := g_compsdb[n]
	if !ok {
		g_compsdb[n] = nil, false
	}
	return ok
}

func (self *appMgr) Configure() StatusCode {
	self.evtproc = NewEvtProcessor("evt-proc")
	//self.evtsel  = 

	return StatusCode(0)
}

func (self *appMgr) Initialize() StatusCode {
	allgood := true
	self.MsgInfo("initialize...\n")

	self.MsgVerbose("components-map: %v\n", g_compsdb)
	svcs_prop, ok := self.GetProperty("Svcs").([]string)
	if ok {
		self.MsgInfo("svcs...\n")
		for _,svc_name := range svcs_prop {
			isvc := self.GetService(svc_name)
			sc := isvc.InitializeSvc()
			if sc != StatusCode(0) {
				self.MsgError("pb initializing [%s] !\n",isvc.CompName())
				allgood = false
			}
		}
		_ = self.evtproc.(IService).InitializeSvc()
		self.MsgInfo("svcs... [done]\n")
	}

	algs_prop, ok := self.GetProperty("Algs").([]string)
	if ok {
		self.MsgInfo("algs...\n")
		for _,alg_name := range algs_prop {
			ialg,isalg := self.GetComp(alg_name).(IAlgorithm)
			if isalg {
				sc := ialg.Initialize()
				if sc != StatusCode(0) {
					self.MsgError("pb initializing [%s] !\n",ialg.CompName())
					allgood = false
				} else {
					self.MsgDebug("correctly initialized [%T/%s]\n",
						ialg, ialg.CompName())
				}
			}
		}
		self.MsgInfo("algs... [done]\n")
	}
	if allgood {
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *appMgr) Start() StatusCode {
	self.MsgInfo("start...\n")
	return StatusCode(0)
}

func (self *appMgr) Run() StatusCode {
	self.MsgInfo("run...\n")
	// init
	sc := self.Initialize()
	if !sc.IsSuccess() {
		return sc
	}

	// start
	sc = self.Start()
	if !sc.IsSuccess() {
		return sc
	}

	// evtloop-run
	sc = self.evtproc.ExecuteRun(10)
	if !sc.IsSuccess() {
		return sc
	}

	// stop
	sc = self.Stop()
	if !sc.IsSuccess() {
		return sc
	}

	// fini
	sc = self.Finalize()
	if !sc.IsSuccess() {
		return sc
	}

	return self.Terminate()
}

func (self *appMgr) Stop() StatusCode {
	self.MsgInfo("stop...\n")
	return StatusCode(0)
}

func (self *appMgr) Finalize() StatusCode {
	self.MsgInfo("finalize...\n")
	allgood := true

	svcs_prop, ok := self.GetProperty("Svcs").([]string)
	if ok {
		self.MsgInfo("svcs...\n")
		for _,svc_name := range svcs_prop {
			isvc := self.GetService(svc_name)
			sc := isvc.FinalizeSvc()
			if sc != StatusCode(0) {
				self.MsgError("pb finalizing [%s] !\n",isvc.CompName())
				allgood = false
			}
		}
		_ = self.evtproc.(IService).FinalizeSvc()
		self.MsgInfo("svcs... [done]\n")
	}

	algs_prop, ok := self.GetProperty("Algs").([]string)
	if ok {
		self.MsgInfo("algs...\n")
		for _,alg_name := range algs_prop {
			ialg,isalg := self.GetComp(alg_name).(IAlgorithm)
			if isalg {
				sc := ialg.Finalize()
				if sc != StatusCode(0) {
					self.MsgError("pb finalizing [%s] !\n",ialg.CompName())
					allgood = false
				} else {
					self.MsgDebug("correctly finalized [%T/%s]\n",
						ialg, ialg.CompName())
				}
			}
		}
		self.MsgInfo("algs... [done]\n")
	}
	if allgood {
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *appMgr) Terminate() StatusCode {
	self.MsgInfo("terminate...\n")
	return StatusCode(0)
}

func NewAppMgr() IAppMgr {
	self := &appMgr{}
	self.properties.props = make(map[string]interface{})
	self.name = "app-mgr"
	self.jobo = "foo.py"
	self.msgstream = msgstream{name:self.name, level:LVL_INFO}

	self.svcMgr.services = make(map[string]IService)
	self.algMgr.algs = make(map[string]IAlgorithm)


	self.mgrs = make(map[string]IComponentMgr)
	self.mgrs["svcmgr"] = &self.svcMgr
	self.mgrs["algmgr"] = &self.algMgr
	
	g_compsdb[self.name] = self

	// completing bootstrap
	g_isvcloc = self

	return self
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
