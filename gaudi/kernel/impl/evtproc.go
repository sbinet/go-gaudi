package kernel

import "time"

type evtProc struct {
	name string
	algs []IAlgorithm
}

func (e *evtProc) CompType() string {
	return "gaudi.kernel.evtProc"
}

func (e *evtProc) CompName() string {
	return e.name
}

func (e *evtProc) ExecuteEvent(ctx IEvtCtx) StatusCode {
	if ctx != nil {
		println(e.name, " executing event...", ctx.(int))
		for i,alg := range e.algs {
			sc := alg.Execute()
			if sc != 0 {
				println(e.name, "pb executing alg #",i,"(",alg.CompName(),")")
				return StatusCode(1)
			}
		}
		return StatusCode(0)
	}
	return StatusCode(-1)
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
	return &evtProc{name, []IAlgorithm{}}
}

// ---

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

/* EOF */
