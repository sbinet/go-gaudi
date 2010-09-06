package kernel

type StatusCode int

type IComponent interface {
	CompName() string
	CompType() string
}

type IService interface {
	Initialize() StatusCode
	Finalize() StatusCode
}

type IAlgorithm interface {
	Initialize() StatusCode
	Execute() StatusCode
	Finalize() StatusCode
}

type IAlgTool interface {
	Initialize() StatusCode
	Finalize() StatusCode
}

type IEvtCtx interface {
	
}

type IEvtProcessor interface {
	ExecuteEvent(evtctx IEvtCtx) StatusCode
	ExecuteRun(maxevt int) StatusCode
	NextEvent(maxevt int) StatusCode
	StopRun() StatusCode
}

type IEvtSelector interface {
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

type IDataStore interface {
	Get(key string) (chan *interface{}, bool)
	Put(key string, value *interface{})
	Has(key string) bool
}
