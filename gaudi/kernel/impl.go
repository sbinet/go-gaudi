package kernel

import "time"

type storeDict map[string]interface{}

type DataStore struct {
	name  string
	store storeDict
}

func NewDataStore(name string) *DataStore {
	store := &DataStore{name:name}
	store.store = make(storeDict)
	return store
}

func (d *DataStore) Put(key string, value *interface{}) bool {
	ok := true
	d.store[key] = value
	return ok
}

func (d *DataStore) Has(key string) bool {
	_,ok := d.store[key]
	return ok
}

func (d *DataStore) Get(key string) (chan *interface{}, bool) {
	out := make(chan *interface{})
	v,ok := d.store[key]
	if ok {
		out <- &v
	} else {
		out <- nil
	}
	return out, ok
}

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

type evtProc struct {
	name string
}

func (e *evtProc) ExecuteEvent(ctx IEvtCtx) StatusCode {
	if ctx != nil {
		println(e.name, " executing event...", ctx.(int))
		return StatusCode(0)
	}
	return StatusCode(-1)
}

func (e *evtProc) test_0() {
	handle := func(queue chan int) StatusCode {
		sc := StatusCode(0)
		for i := range queue {
			println("   --> handling [",i,"]...")
			sc = e.ExecuteEvent(i)
		}
		return sc
	}

	max_in_flight := 4
	serve := func(queue chan int, quit chan bool) StatusCode {
		for i := 0; i < max_in_flight; i++ {
			go handle(queue)
		}
		<-quit // wait to be told to exit
		return StatusCode(0)
	}

	quit := make(chan bool)

	println("-- filling the event queue...")
	queue := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			queue <- i
		}
	}()
	println("-- starting to serve 20 events...")
	go serve(queue, quit)
	println("-- requests sent...")
	time.Sleep(2000000000)
	quit <- true
	println("-- done.")
}

func (e *evtProc) ExecuteRun(evtmax int) StatusCode {
	println(e.name, "execute-run [", evtmax, "]...")
	sc := e.NextEvent(evtmax)
	return sc
}

type evtRequest struct {
	idx  int
	sc   StatusCode
}

func (e *evtProc) NextEvent(evtmax int) StatusCode {

	handle := func(evt *evtRequest, out_queue chan *evtRequest) {
		println(e.name, " nextEvent[", evt.idx, "]...")
		evt.sc = e.ExecuteEvent(evt.idx)
		out_queue <- evt
	}

	serve_evts := func(in_evt_queue, out_evt_queue chan *evtRequest, quit chan bool) {
		for {
			select {
			case ievt := <-in_evt_queue:
				go handle(ievt, out_evt_queue)
			case <-quit:
				//println("quit requested !")
				return
			}
		}
	}

	start_evt_server := func(nworkers int) (in_evt_queue, 
		                                    out_evt_queue chan *evtRequest,
		                                    quit chan bool) {
		in_evt_queue = make(chan *evtRequest, nworkers)
		out_evt_queue = make(chan *evtRequest)
		quit = make(chan bool)
		go serve_evts(in_evt_queue, out_evt_queue, quit)
		return in_evt_queue, out_evt_queue, quit
	}

	const nworkers = 4
	in_evt_queue, out_evt_queue, quit := start_evt_server(nworkers)
	//println(e.name, "sending requests...")
	for i:=0; i<evtmax; i++ {
		in_evt_queue <- &evtRequest{i, StatusCode(0)}
	}
	//println(e.name, "sending requests... [done]")
	n_fails := 0
	n_processed := 0
	for evt := range out_evt_queue {
		//println(e.name, "out-evt-queue:",evt.idx, evt.sc)
		if evt.sc != StatusCode(0) {
			n_fails++
		}
		n_processed++
		if n_processed == evtmax {
			quit <- true
			close(out_evt_queue)
			//println("closing evt server...")
			break
		}
	}
	if n_fails != 0 {
		return StatusCode(1)
	}
	return StatusCode(0)
}

func (e *evtProc) StopRun() StatusCode {
	println(e.name, " stopping run...")
	return StatusCode(0)
}

func NewEvtProcessor(name string) IEvtProcessor {
	return &evtProc{name}
}
