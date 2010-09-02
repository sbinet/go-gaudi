package main

import "gaudi/kernel"
import "fmt"
import "time"

type AppMgr struct {
	name string
	jobo string
	
	evtproc kernel.IEvtProcessor
	evtsel  kernel.IEvtSelector
}

func (app *AppMgr) CompType() string {
	return "gaudi.AppMgr"
}

func (app *AppMgr) CompName() string {
	return app.name
}

func (app *AppMgr) ExecuteRun() kernel.StatusCode {
	println(app.name, " execute-run...")
	sc := app.evtproc.ExecuteRun(10)
	return sc
}
func (app *AppMgr) Configure() kernel.StatusCode {
	app.evtproc = NewEvtProcessor("evt-proc")
	//app.evtsel  = 

	return kernel.StatusCode(-1)
}

type evtProc struct {
	name string
}

func (e *evtProc) ExecuteEvent(ctx kernel.IEvtCtx) kernel.StatusCode {
	if ctx != nil {
		println(e.name, " executing event...", ctx.(int))
		return kernel.StatusCode(0)
	}
	return kernel.StatusCode(-1)
}

func (e *evtProc) test_0() {
	handle := func(queue chan int) kernel.StatusCode {
		sc := kernel.StatusCode(0)
		for i := range queue {
			println("   --> handling [",i,"]...")
			sc = e.ExecuteEvent(i)
		}
		return sc
	}

	max_in_flight := 4
	serve := func(queue chan int, quit chan bool) kernel.StatusCode {
		for i := 0; i < max_in_flight; i++ {
			go handle(queue)
		}
		<-quit // wait to be told to exit
		return kernel.StatusCode(0)
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

func (e *evtProc) ExecuteRun(evtmax int) kernel.StatusCode {
	println(e.name, "execute-run [", evtmax, "]...")
	sc := e.NextEvent(evtmax)
	return sc
}

type evtRequest struct {
	idx  int
	sc   kernel.StatusCode
}

func (e *evtProc) NextEvent(evtmax int) kernel.StatusCode {

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
		in_evt_queue <- &evtRequest{i, kernel.StatusCode(0)}
	}
	//println(e.name, "sending requests... [done]")
	n_fails := 0
	n_processed := 0
	for evt := range out_evt_queue {
		//println(e.name, "out-evt-queue:",evt.idx, evt.sc)
		if evt.sc != kernel.StatusCode(0) {
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
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (e *evtProc) StopRun() kernel.StatusCode {
	println(e.name, " stopping run...")
	return kernel.StatusCode(0)
}

func NewEvtProcessor(name string) kernel.IEvtProcessor {
	return &evtProc{name}
}

func handle_icomponent(c kernel.IComponent) {
	fmt.Printf(":: handle_icomponent(%s)...\n", c.CompName())
}

func main() {
	fmt.Print("::: gaudi\n")
	app := &AppMgr{name:"app-mgr", jobo:"foo.py"}

	fmt.Printf(" -> created [%s/%s]\n", app.CompType(), app.CompName())

	handle_icomponent(app)
	fmt.Printf("%s\n", app)

	println("::: configure...")
	sc := app.Configure()
	println("::: configure... [", sc, "]")
	println("::: run...")
	sc = app.ExecuteRun()
	println("::: run... [", sc, "]")

	// println("::: testing event server...")
	// e := app.evtproc.(*evtProc)
	// e.test_0()
	// println("::: testing event server... [done]")
	fmt.Print("::: bye.\n")
}
