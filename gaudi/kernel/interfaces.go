package kernel

import "fmt"

type comps_db map[string]IComponent

/// the main entry point to the gaudi framework: the service locator
var g_isvcloc ISvcLocator = nil
/// the central repository of all gaudi components
var g_compsdb comps_db

type StatusCode int

func (sc StatusCode) String() string {
	return fmt.Sprintf("%i", int(sc))
}

func (sc StatusCode) IsSuccess() bool {
	return sc == StatusCode(0)
}

func (sc StatusCode) IsFailure() bool {
	return !sc.IsSuccess()
}

func (sc StatusCode) IsRecoverable() bool {
	return sc == StatusCode(2)
}

func GetSvcLocator() ISvcLocator {
	return g_isvcloc
}

func ComponentMgr() IComponentMgr {
	isvcloc := GetSvcLocator()
	imgr, ok := isvcloc.(IComponentMgr)
	if ok {
		return imgr
	}
	return nil
}

func RegisterComp(c IComponent) bool {
	if c == nil {
		return false
	}
	n := c.CompName()
	oldcomp, already_there := g_compsdb[n]
	if already_there {
		if oldcomp == c {
			// double registration of the same component...
			// silly but harmless.
			return true
		}
		// already existing component with that same name !
		err := fmt.Sprintf("a component with name [%s] was already registered ! (old-type: %T, new-type: %T)",
			n, oldcomp, c);
		panic(err)
	}
	//fmt.Printf("--> registering [%T/%s]...\n", c, n)
	g_compsdb[n] = c
	//fmt.Printf("--> registering [%T/%s]... [done]\n", c, n)
	return true
}
type IComponent interface {
	CompName() string
	CompType() string
}

type IComponentMgr interface {
	GetComp(n string) IComponent
	GetComps() []IComponent
	HasComp(n string) bool
}

type Property struct {
	Name string
	Value interface{}
}
type IProperty interface {
	/// set the property value
	SetProperty(name string, value interface{}) StatusCode
	/// get the property value by name
	GetProperty(name string) interface{}
	/// get the list of properties
	GetProperties() []Property
}

type IService interface {
	IComponent
	InitializeSvc() StatusCode
	FinalizeSvc() StatusCode
}

type IAlgorithm interface {
	IComponent
	Initialize() StatusCode
	Execute(evtctx IEvtCtx) StatusCode
	Finalize() StatusCode
}

type IAlgTool interface {
	IComponent
	InitializeTool() StatusCode
	FinalizeTool() StatusCode
}

type IEvtCtx interface {
	
}

type IEvtProcessor interface {
	IComponent
	ExecuteEvent(evtctx IEvtCtx) StatusCode
	ExecuteRun(maxevt int) StatusCode
	NextEvent(maxevt int) StatusCode
	StopRun() StatusCode
}

type IEvtSelector interface {
	IComponent
	CreateContext(ctx *IEvtCtx) StatusCode
	Next(ctx *IEvtCtx, jump int) StatusCode
	Previous(ctx *IEvtCtx, jump int) StatusCode
	Last(ctx *IEvtCtx) StatusCode
	Rewind(ctx *IEvtCtx) StatusCode
}

type IAppMgr interface {
	IComponent
	Configure() StatusCode
	Initialize() StatusCode
	Start() StatusCode
	/// Run the complete job (from Initialize to Terminate)
	Run() StatusCode
	Stop() StatusCode
	Finalize() StatusCode
	Terminate() StatusCode
}

type IAlgMgr interface {
	//IComponent
	AddAlgorithm(alg IAlgorithm) StatusCode
	RemoveAlgorithm(alg IAlgorithm) StatusCode
	HasAlgorithm(algname string) bool
}

type ISvcMgr interface {
	//IComponent
	AddService(svc string) StatusCode
	RemoveService(svc string) StatusCode
	HasService(svc string) StatusCode
}

type ISvcLocator interface {
	//IComponent
	GetService(svc string) IService
	GetServices() []IService
	ExistsService(svc string) bool
}

type IDataStore interface {
	IComponent
	Get(key string) (chan *interface{}, bool)
	Put(key string, value *interface{}) bool
	Has(key string) bool
}
