package kernel

import "time"
//import "fmt"

type evtproc struct {
	Service
	algs []IAlgorithm
}

func (self *evtproc) InitializeSvc() StatusCode {

	sc := self.Service.InitializeSvc()
	if !sc.IsSuccess() {
		return sc
	}

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
	if ictx != nil {
		ctx := ictx.(int)
		{
			self.MsgInfo("-- clear evt-store [ctx:%03v]...\n", ctx)
			smgr := GetSvcLocator().GetService("evt-store").(IDataStoreMgr)
			store := smgr.Store(ictx).(IDataStoreClearer)
			if !store.ClearStore().IsSuccess() {
				self.MsgWarning("could not clear store [evt-store]\n")
			}
			self.MsgInfo("-- clear evt-store [ctx:%03v]... [done]\n", ctx)
		}
		self.MsgInfo("executing event [%v]... (#algs: %v)\n", ctx, len(self.algs))
		for i,alg := range self.algs {
			self.MsgInfo("-- ctx:%03v --> [%s]...\n", ctx, alg.CompName())
			sc := alg.Execute(ictx)
			if sc != 0 {
				self.MsgError("pb executing alg #%v (%s) for ctx:%v\n",
					i,alg.CompName(), ictx)
				return StatusCode(1)
			}
		}
		return StatusCode(0)
	}
	return StatusCode(-1)
}

func (self *evtproc) ExecuteRun(evtmax int) StatusCode {
	self.MsgInfo("execute-run [%v]\n", evtmax)
	sc := self.NextEvent(evtmax)
	return sc
}

type evtrequest struct {
	idx  int
	sc   StatusCode
}

func (self *evtproc) NextEvent(evtmax int) StatusCode {

	handle := func(evt *evtrequest, out_queue chan *evtrequest) {
		self.MsgInfo("nextEvent[%v]...\n", evt.idx)
		evt.sc = self.ExecuteEvent(evt.idx)
		out_queue <- evt
	}

	serve_evts := func(in_evt_queue, out_evt_queue chan *evtrequest, quit chan bool) {
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
		                                    out_evt_queue chan *evtrequest,
		                                    quit chan bool) {
		in_evt_queue = make(chan *evtrequest, nworkers)
		out_evt_queue = make(chan *evtrequest)
		quit = make(chan bool)
		go serve_evts(in_evt_queue, out_evt_queue, quit)
		return in_evt_queue, out_evt_queue, quit
	}

	const nworkers = 1
	smgr := GetSvcLocator().GetService("evt-store").(IDataStoreMgr)
	if !smgr.SetNbrStreams(nworkers).IsSuccess() {
		self.MsgWarning("problem setting the correct number of parallel evt-stores\n")
		return StatusCode(1)
	}

	in_evt_queue, out_evt_queue, quit := start_evt_server(nworkers)
	//println(e.CompName(), "sending requests...")
	for i:=0; i<evtmax; i++ {
		in_evt_queue <- &evtrequest{i, StatusCode(0)}
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
	return self
}

// ---

func (self *evtproc) test_0() {
	handle := func(queue chan int) StatusCode {
		sc := StatusCode(0)
		for i := range queue {
			self.MsgInfo("   --> handling [%i]...\n",i)
			sc = self.ExecuteEvent(i)
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
var _ = IComponent(&evtproc{})
var _ = IEvtProcessor(&evtproc{})
var _ = IProperty(&evtproc{})
var _ = IService(&evtproc{})

/* EOF */
