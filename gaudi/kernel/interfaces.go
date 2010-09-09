package kernel

type StatusCode int

type IComponent interface {
	CompName() string
	CompType() string
}

type IComponentMgr interface {
	GetComp(n string) *IComponent
	GetComps() []*IComponent
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
	Execute() StatusCode
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
	IComponent
	AddAlgorithm(alg IAlgorithm) StatusCode
	RemoveAlgorithm(alg IAlgorithm) StatusCode
	HasAlgorithm(algname string) bool
}

type ISvcMgr interface {
	IComponent
	AddService(svc string) StatusCode
	RemoveService(svc string) StatusCode
	HasService(svc string) StatusCode
}

type IDataStore interface {
	IComponent
	Get(key string) (chan *interface{}, bool)
	Put(key string, value *interface{})
	Has(key string) bool
}
