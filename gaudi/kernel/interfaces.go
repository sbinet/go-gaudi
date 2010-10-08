package kernel

import "fmt"
import "os"

type comps_db map[string]IComponent

/// the main entry point to the gaudi framework: the service locator
var g_isvcloc ISvcLocator = nil
/// the central repository of all gaudi components
var g_compsdb comps_db

type Error interface {
	Code() int
	Err() os.Error
	IsSuccess() bool
	IsFailure() bool
	IsRecoverable() bool
}

type statuscode struct {
	code int
	err  os.Error
}

func StatusCode(code int) Error {
	return &statuscode{code:code, err:nil}
}

func (sc *statuscode) Code() int {
	return sc.code
}

func (sc *statuscode) Err() os.Error {
	return sc.err
}

func (sc *statuscode) String() string {
	return fmt.Sprintf("code:%i err:%v", sc.code, sc.err)
}

func (sc *statuscode) IsSuccess() bool {
	return sc.code == 0
}

func (sc *statuscode) IsFailure() bool {
	return !sc.IsSuccess()
}

func (sc *statuscode) IsRecoverable() bool {
	return sc.code == 2
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
	/// declare a property by name and default value
	DeclareProperty(name string, value interface{})
	/// set the property value
	SetProperty(name string, value interface{}) Error
	/// get the property value by name
	GetProperty(name string) interface{}
	/// get the list of properties
	GetProperties() []Property
}

type IService interface {
	IComponent
	//SysInitializeSvc() Error
	//SysFinalizeSvc() Error

	InitializeSvc() Error
	FinalizeSvc() Error
}

type IAlgorithm interface {
	IComponent
	//SysInitialize() Error
	//SysExecute(evtctx IEvtCtx) Error
	//SysFinalize() Error
	Initialize() Error
	Execute(evtctx IEvtCtx) Error
	Finalize() Error
}

type IAlgTool interface {
	IComponent
	//SysInitializeTool() Error
	//SysFinalizeTool() Error

	InitializeTool() Error
	FinalizeTool() Error
}

type DataStore map[string]interface{}
type IEvtCtx interface {
	Idx() int
	Store() *DataStore
	//Id() int
}

type IEvtProcessor interface {
	IComponent
	ExecuteEvent(evtctx IEvtCtx) Error
	ExecuteRun(maxevt int) Error
	NextEvent(maxevt int) Error
	StopRun() Error
}

type IEvtSelector interface {
	IComponent
	CreateContext(ctx *IEvtCtx) Error
	Next(ctx *IEvtCtx, jump int) Error
	Previous(ctx *IEvtCtx, jump int) Error
	Last(ctx *IEvtCtx) Error
	Rewind(ctx *IEvtCtx) Error
}

type IAppMgr interface {
	IComponent
	Configure() Error
	Initialize() Error
	Start() Error
	/// Run the complete job (from Initialize to Terminate)
	Run() Error
	Stop() Error
	Finalize() Error
	Terminate() Error
}

type IAlgMgr interface {
	//IComponent
	AddAlgorithm(alg IAlgorithm) Error
	RemoveAlgorithm(alg IAlgorithm) Error
	HasAlgorithm(algname string) bool
}

type ISvcMgr interface {
	//IComponent
	AddService(svc string) Error
	RemoveService(svc string) Error
	HasService(svc string) Error
}

type ISvcLocator interface {
	//IComponent
	GetService(svc string) IService
	GetServices() []IService
	ExistsService(svc string) bool
}

type IDataStore interface {
	//IComponent
	Get(key string) interface{}
	Put(key string, value interface{})
	Has(key string) bool
	//Keys() []string // ??
}

type IDataStoreMgr interface {
	IComponent
	Store(ctx IEvtCtx) IDataStore
}

type IMessager interface {
	Msg(lvl OutputLevel, format string, a ...interface{}) (int, os.Error)
	MsgVerbose(format string, a ...interface{}) (int, os.Error)
	MsgDebug(format string, a ...interface{}) (int, os.Error)
	MsgInfo(format string, a ...interface{}) (int, os.Error)
	MsgWarning(format string, a ...interface{}) (int, os.Error)
	MsgError(format string, a ...interface{}) (int, os.Error)
	MsgFatal(format string, a ...interface{}) (int, os.Error)
	MsgAlways(format string, a ...interface{}) (int, os.Error)
}
/* EOF */
