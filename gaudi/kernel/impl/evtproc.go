package kernel

import "time"
//import "fmt"

type evtProc struct {
	Service
	algs []IAlgorithm
}

func (e *evtProc) InitializeSvc() StatusCode {
	sc := e.Service.InitializeSvc()
	if !sc.IsSuccess() {
		return sc
	}

	svcloc := GetSvcLocator()
	if svcloc == nil {
		e.MsgError("could not retrieve ISvcLocator !\n")
		return StatusCode(1)
	}
	appmgr := svcloc.(IComponentMgr).GetComp("app-mgr")
	if appmgr == nil {
		e.MsgError("could not retrieve 'app-mgr'\n")
	}
	propmgr := appmgr.(IProperty)
	alg_names := propmgr.GetProperty("Algs").([]string)
	e.MsgInfo("got alg-names: %v\n", alg_names)

	if len(alg_names)>0 {
		comp_mgr := appmgr.(IComponentMgr)
		e.algs = make([]IAlgorithm, len(alg_names))
		for i,alg_name := range alg_names {
			ialg,isalg := comp_mgr.GetComp(alg_name).(IAlgorithm)
			if isalg {
				e.algs[i] = ialg
			}
		}
	}
	e.MsgInfo("got alg-list: %v\n", e.algs)

	return StatusCode(0)
}

func (e *evtProc) ExecuteEvent(ictx IEvtCtx) StatusCode {
	if ictx != nil {
		ctx := ictx.(int)
		e.MsgInfo("executing event [%v]...\n", ctx)
		for i,alg := range e.algs {
			e.MsgInfo("-- ctx:%03v --> [%s]...\n", ctx, alg.CompName())
			sc := alg.Execute(ictx)
			if sc != 0 {
				e.MsgError("pb executing alg #%v (%s) for ctx:%v\n",
					i,alg.CompName(), ictx)
				return StatusCode(1)
			}
		}
		return StatusCode(0)
	}
	return StatusCode(-1)
}

func (e *evtProc) ExecuteRun(evtmax int) StatusCode {
	e.MsgInfo("execute-run [%v]\n", evtmax)
	sc := e.NextEvent(evtmax)
	return sc
}

type evtRequest struct {
	idx  int
	sc   StatusCode
}

func (e *evtProc) NextEvent(evtmax int) StatusCode {

	handle := func(evt *evtRequest, out_queue chan *evtRequest) {
		e.MsgInfo("nextEvent[%v]...\n", evt.idx)
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
	//println(e.CompName(), "sending requests...")
	for i:=0; i<evtmax; i++ {
		in_evt_queue <- &evtRequest{i, StatusCode(0)}
	}
	//println(e.CompName(), "sending requests... [done]")
	n_fails := 0
	n_processed := 0
	for evt := range out_evt_queue {
		//println(e.CompName(), "out-evt-queue:",evt.idx, evt.sc)
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
	e.MsgInfo("stopping run...\n")
	return StatusCode(0)
}

func NewEvtProcessor(name string) IEvtProcessor {
	p := &evtProc{}
	
	//p.properties.props = make(map[string]interface{})
	//p.name = name
	p.algs = []IAlgorithm{}
	_ = NewSvc(&p.Service, "evtProc", name)
	return p
}

// ---

func (e *evtProc) test_0() {
	handle := func(queue chan int) StatusCode {
		sc := StatusCode(0)
		for i := range queue {
			e.MsgInfo("   --> handling [%i]...\n",i)
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

	e.MsgInfo("-- filling the event queue...\n")
	queue := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			queue <- i
		}
	}()
	e.MsgInfo("-- starting to serve 20 events...\n")
	go serve(queue, quit)
	e.MsgInfo("-- requests sent...\n")
	time.Sleep(2000000000)
	quit <- true
	e.MsgInfo("-- done.\n")
}

// check implementations match interfaces
var _ = IComponent(&evtProc{})
var _ = IEvtProcessor(&evtProc{})
var _ = IProperty(&evtProc{})
var _ = IService(&evtProc{})

/* EOF */
