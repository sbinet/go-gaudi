package kernel

import "time"
//import "fmt"

// --- evt state ---
type evtstate struct {
	idx  int
	sc   StatusCode
	data DataStore
}

func new_evtstate(idx int) *evtstate {
	self := &evtstate{idx:idx, sc:StatusCode(0), data:make(DataStore)}
	self.data["evt-store"] = make(DataStore)
	return self
}

func (self *evtstate) Idx() int {
	return self.idx
}
func (self *evtstate) Store() *DataStore {
	return &self.data
}

// --- evt processor ---
type evtproc struct {
	Service
	algs []IAlgorithm
	nworkers int
}

func (self *evtproc) InitializeSvc() StatusCode {

	sc := self.Service.InitializeSvc()
	if !sc.IsSuccess() {
		return sc
	}
	self.nworkers = self.GetProperty("NbrWorkers").(int)
	//self.nworkers = 1
	self.MsgInfo("n-workers: %v\n", self.nworkers)
	svcloc := GetSvcLocator()
	if svcloc == nil {
		self.MsgError("could not retrieve ISvcLocator !\n")
		return StatusCode(1)
	}
	appmgr := svcloc.(IComponentMgr).GetComp("app-mgr")
	if appmgr == nil {
		self.MsgError("could not retrieve 'app-mgr'\n")
	}
	propmgr := appmgr.(IProperty)
	alg_names := propmgr.GetProperty("Algs").([]string)
	self.MsgInfo("got alg-names: %v\n", alg_names)

	if len(alg_names)>0 {
		comp_mgr := appmgr.(IComponentMgr)
		self.algs = make([]IAlgorithm, len(alg_names))
		for i,alg_name := range alg_names {
			ialg,isalg := comp_mgr.GetComp(alg_name).(IAlgorithm)
			if isalg {
				self.algs[i] = ialg
			}
		}
	}
	self.MsgInfo("got alg-list: %v\n", self.algs)

	return StatusCode(0)
}

func (self *evtproc) ExecuteEvent(ictx IEvtCtx) StatusCode {
	ctx := ictx.Idx()
	self.MsgInfo("executing event [%v]... (#algs: %v)\n", ctx, len(self.algs))
	for i,alg := range self.algs {
		self.MsgDebug("-- ctx:%03v --> [%s]...\n", ctx, alg.CompName())
		if !alg.Execute(ictx).IsSuccess() {
			self.MsgError("pb executing alg #%v (%s) for ctx:%v\n",
				i,alg.CompName(), ictx.Idx())
			return StatusCode(1)
		}
	}
	self.MsgInfo("data: %v\n",*ictx.Store())
	return StatusCode(0)
}

func (self *evtproc) ExecuteRun(evtmax int) StatusCode {
	self.MsgInfo("execute-run [%v]\n", evtmax)
	sc := self.NextEvent(evtmax)
	return sc
}

func (self *evtproc) NextEvent(evtmax int) StatusCode {

	if self.nworkers > 1 {
		return self.mp_NextEvent(evtmax)
	}
	return self.seq_NextEvent(evtmax)
}

func (self *evtproc) seq_NextEvent(evtmax int) StatusCode {
	
	self.MsgInfo("nextEvent[%v]...\n", evtmax)
	for i:=0; i<evtmax; i++ {
		ctx := new_evtstate(i)
		if !self.ExecuteEvent(ctx).IsSuccess() {
			self.MsgError("failed to execute evt idx %03v\n", i)
			return StatusCode(1)
		}
	}
	return StatusCode(0)
}
func (self *evtproc) mp_NextEvent(evtmax int) StatusCode {

	handle := func(evt *evtstate, out_queue chan *evtstate) {
		self.MsgInfo("nextEvent[%v]...\n", evt.idx)
		evt.sc = self.ExecuteEvent(evt)
		out_queue <- evt
	}

	serve_evts := func(in_evt_queue, out_evt_queue chan *evtstate, quit chan bool) {
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
		                                    out_evt_queue chan *evtstate,
		                                    quit chan bool) {
		in_evt_queue = make(chan *evtstate, nworkers)
		out_evt_queue = make(chan *evtstate)
		quit = make(chan bool)
		go serve_evts(in_evt_queue, out_evt_queue, quit)
		return in_evt_queue, out_evt_queue, quit
	}

	/*
	smgr := GetSvcLocator().GetService("evt-store").(IDataStoreMgr)
	if !smgr.SetNbrStreams(self.nworkers).IsSuccess() {
		self.MsgWarning("problem setting the correct number of parallel evt-stores\n")
		return StatusCode(1)
	}
	 */

	in_evt_queue, out_evt_queue, quit := start_evt_server(self.nworkers)
	//println(e.CompName(), "sending requests...")
	for i:=0; i<evtmax; i++ {
		in_evt_queue <- new_evtstate(i)
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

func (self *evtproc) StopRun() StatusCode {
	self.MsgInfo("stopping run...\n")
	return StatusCode(0)
}

func NewEvtProcessor(name string) IEvtProcessor {
	self := &evtproc{}
	
	//self.properties.props = make(map[string]interface{})
	//self.name = name
	self.algs = []IAlgorithm{}
	_ = NewSvc(&self.Service, "kernel.evtproc", name)
	RegisterComp(self)
	self.DeclareProperty("NbrWorkers", 4)
	return self
}

// ---

func (self *evtproc) test_0() {
	handle := func(queue chan int) StatusCode {
		sc := StatusCode(0)
		for i := range queue {
			ctx := new_evtstate(i)
			self.MsgInfo("   --> handling [%i]...\n",i)
			sc = self.ExecuteEvent(ctx)
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

	self.MsgInfo("-- filling the event queue...\n")
	queue := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			queue <- i
		}
	}()
	self.MsgInfo("-- starting to serve 20 events...\n")
	go serve(queue, quit)
	self.MsgInfo("-- requests sent...\n")
	time.Sleep(2000000000)
	quit <- true
	self.MsgInfo("-- done.\n")
}

// check implementations match interfaces
var _ = IEvtCtx(&evtstate{})

var _ = IComponent(&evtproc{})
var _ = IEvtProcessor(&evtproc{})
var _ = IProperty(&evtproc{})
var _ = IService(&evtproc{})

/* EOF */
