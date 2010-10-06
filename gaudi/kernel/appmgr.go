package kernel

//import "fmt"

///////////////////////////////////////////////////////////////////////////////
// svc-mgr

type svcmgr struct {
	services map[string]IService
}

func (self *svcmgr) GetComp(n string) IComponent {
	c,ok := self.services[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (self *svcmgr) GetComps() []IComponent {
	comps := make([]IComponent, len(self.services))
	i := 0
	for _,v := range self.services {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (self *svcmgr) HasComp(n string) bool {
	_,ok := self.services[n]
	if !ok {
		self.services[n] = nil, false
	}
	return ok
}

func (self *svcmgr) AddService(svc string) StatusCode {
	isvc,ok := g_compsdb[svc].(IService)
	if !ok {
		//fmt.Printf("** AddService(%s) FAILED !\n", svc)
		return StatusCode(1)
	}
	self.services[svc] = isvc
	return StatusCode(0)
}

func (self *svcmgr) RemoveService(svc string) StatusCode {
	if self.HasService(svc) == StatusCode(0) {
		self.services[svc] = nil, false
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *svcmgr) HasService(svc string) StatusCode {
	if self.HasComp(svc) {
		//fmt.Printf(":: HasService(%s) - true\n", svc)
		return StatusCode(0)
	}
	return StatusCode(1)
}

func (self *svcmgr) GetService(svc string) IService {
	if self.HasService(svc).IsSuccess() {
		//fmt.Printf("-- GetService(%s)...\n", svc)
		isvc := self.services[svc]
		//fmt.Printf("-- GetService(%s)... [done]\n", svc)
		return isvc
	}
	return nil
}

func (self *svcmgr) GetServices() []IService {
	svcs := make([]IService, len(self.services))
	i := 0
	for _,v := range self.services {
		svcs[i] = v
		i++
	}
	return svcs
}

func (self *svcmgr) ExistsService(svc string) bool {
	return self.HasService(svc) == StatusCode(0)
}


//////////////////////////////////////////////////////////////////////////////
// alg-mgr

type algmgr struct {
	algs map[string]IAlgorithm
}

func (self *algmgr) GetComp(n string) IComponent {
	c,ok := self.algs[n]
	if ok {
		return c.(IComponent)
	}
	return nil
}

func (self *algmgr) GetComps() []IComponent {
	comps := make([]IComponent, len(self.algs))
	i := 0
	for _,v := range self.algs {
		comps[i] = v.(IComponent)
		i++
	}
	return comps
}

func (self *algmgr) HasComp(n string) bool {
	_,ok := self.algs[n]
	if !ok {
		self.algs[n] = nil, false
	}
	return ok
}

func (self *algmgr) AddAlgorithm(alg IAlgorithm) StatusCode {
	self.algs[alg.CompName()] = alg
	return StatusCode(0)
}

func (self *algmgr) RemoveAlgorithm(alg IAlgorithm) StatusCode {
	n := alg.CompName()
	if !self.HasComp(n) {
		return StatusCode(1)
	}
	self.algs[n] = nil, false
	return StatusCode(0)
}

func (self *algmgr) HasAlgorithm(algname string) bool {
	return self.HasComp(algname)
}

//////////////////////////////////////////////////////////////////////////////
// app-mgr

type appmgr struct {
	properties
	msgstream
	svcmgr
	algmgr

	name string
	jobo string
	
	evtproc IEvtProcessor
	evtsel  IEvtSelector

	mgrs    map[string]IComponentMgr
}

func (self *appmgr) CompType() string {
	return "gaudi.kernel.appmgr"
}

func (self *appmgr) CompName() string {
	return self.name
}

func (self *appmgr) GetComp(n string) IComponent {
	if !self.HasComp(n) {
		return nil
	}
	return g_compsdb[n]
}

func (self *appmgr) GetComps() []IComponent {
	comps := make([]IComponent, len(g_compsdb))
	i := 0
	for _,v := range g_compsdb {
		comps[i] = v
		i++
	}
	return comps
}

func (self *appmgr) HasComp(n string) bool {
	_,ok := g_compsdb[n]
	if !ok {
		g_compsdb[n] = nil, false
	}
	return ok
}

func (self *appmgr) Configure() StatusCode {
	//self.evtproc = NewEvtProcessor("evt-proc")
	self.evtproc = self.GetService("evt-proc").(IEvtProcessor)
	//self.evtsel  = 

	return StatusCode(0)
}

func (self *appmgr) Initialize() StatusCode {
	allgood := true
	self.MsgInfo("initialize...\n")

	self.MsgVerbose("components-map: %v\n", g_compsdb)
	svcs_prop, ok := self.GetProperty("Svcs").([]string)
	if ok {
		self.MsgInfo("svcs...\n")
		for _,svc_name := range svcs_prop {
			isvc := self.GetService(svc_name)
			if !isvc.InitializeSvc().IsSuccess() {
				self.MsgError("pb initializing [%s] !\n",isvc.CompName())
				allgood = false
			}
		}
		//_ = self.evtproc.(IService).InitializeSvc()
		self.MsgInfo("svcs... [done]\n")
	}

	algs_prop, ok := self.GetProperty("Algs").([]string)
	if ok {
		self.MsgInfo("algs...\n")
		for _,alg_name := range algs_prop {
			ialg,isalg := self.GetComp(alg_name).(IAlgorithm)
			if isalg {
				if !ialg.Initialize().IsSuccess() {
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

func (self *appmgr) Start() StatusCode {
	self.MsgInfo("start...\n")
	return StatusCode(0)
}

func (self *appmgr) Run() StatusCode {
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

func (self *appmgr) Stop() StatusCode {
	self.MsgInfo("stop...\n")
	return StatusCode(0)
}

func (self *appmgr) Finalize() StatusCode {
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

func (self *appmgr) Terminate() StatusCode {
	self.MsgInfo("terminate...\n")
	return StatusCode(0)
}

func NewAppMgr() IAppMgr {
	self := &appmgr{}
	self.properties.props = make(map[string]interface{})
	self.name = "app-mgr"
	self.jobo = "foo.py"
	self.msgstream = msgstream{name:self.name, level:LVL_INFO}

	self.svcmgr.services = make(map[string]IService)
	self.algmgr.algs = make(map[string]IAlgorithm)


	self.mgrs = make(map[string]IComponentMgr)
	self.mgrs["svcmgr"] = &self.svcmgr
	self.mgrs["algmgr"] = &self.algmgr
	
	g_compsdb[self.name] = self

	// completing bootstrap
	g_isvcloc = self

	return self
}

// check implementations match interfaces
var _ = IAlgMgr(&algmgr{})
var _ = IComponentMgr(&algmgr{})

var _ = IComponentMgr(&svcmgr{})
var _ = ISvcMgr(&svcmgr{})
var _ = ISvcLocator(&svcmgr{})

var _ = IComponent(&appmgr{})
var _ = IComponentMgr(&appmgr{})
var _ = IAlgMgr(&appmgr{})
var _ = ISvcMgr(&appmgr{})
var _ = ISvcLocator(&appmgr{})
var _ = IAppMgr(&appmgr{})
var _ = IProperty(&appmgr{})


/* EOF */