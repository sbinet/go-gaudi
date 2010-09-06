package kernel

type appMgr struct {
	name string
	jobo string
	
	evtproc IEvtProcessor
	evtsel  IEvtSelector
}

func (app *appMgr) CompType() string {
	return "gaudi.appMgr"
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
	return &appMgr{name:"app-mgr", jobo:"foo.py"}
}

