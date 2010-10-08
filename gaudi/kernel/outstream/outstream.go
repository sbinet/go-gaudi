// package containing components to implement gaudi output stream algorithms
package outstream

//import "io"
import "os"
import "gob"
import "json"

import "gaudi/kernel"

var g_keys = []string{
	"njets",
	"ptjets",
	"cnt",
}

// simple interface to gather all encoders
type iwriter interface {
	Encode(v interface{}) os.Error
}

func data_sink(w iwriter) (datachan, chan bool) {
	in   := make(datachan)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case data := <-in:
				err := w.Encode(data)
				if err != nil {
					println("** error **", err)
				}
			case <-quit:
				return
			}
		}
	}()
	return in, quit
}


// --- gob_outstream ---
type gob_outstream struct {
	kernel.Algorithm
	w *os.File
	enc *gob.Encoder
	item_names []string
}

func (self *gob_outstream) Initialize() kernel.StatusCode {
	self.MsgDebug("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	self.item_names = self.GetProperty("Items").([]string)
	fname := self.GetProperty("Output").(string)

	self.MsgInfo("output file: [%v]\n", fname)
	self.MsgInfo("items: %v\n", self.item_names)

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := uint32(0666)
	w,err := os.Open(fname, flag, perm)
	if err != nil {
		self.MsgError("problem while opening file [%v]: %v\n", fname, err)
		return kernel.StatusCode(1)
	}
	self.w = w
	self.enc = gob.NewEncoder(self.w)
	return kernel.StatusCode(0)
}

func (self *gob_outstream) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgDebug("== execute ==\n")
	store := self.EvtStore(ctx)
	if store == nil {
		self.MsgError("could not retrieve evt-store\n")
	}
	var err os.Error
	allgood := true
	for _,k := range self.item_names {
		err = self.enc.Encode(store.Get(k))
		if err != nil {
			self.MsgError("error while writing store content at [%v]: %v\n", 
				k,err)
			allgood = false
		}
	}

	/*
	hdr_offset := 1
	val := make([]interface{}, len(self.item_names)+hdr_offset)
	val[0] = ctx.Idx()

	for i,k := range self.item_names {
		val[i+hdr_offset] = store.Get(k)
	}
	for idx,v := range val {
		err = self.enc.Encode(val)
		if err != nil {
			self.MsgError("error while encoding data [%v|%v]: %v\n", 
				idx, v, err)
			allgood = false
		}
	}
	 */

	if !allgood {
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *gob_outstream) Finalize() kernel.StatusCode {
	self.MsgDebug("== finalize ==\n")
	self.w.Close()

	return kernel.StatusCode(0)
}

// ---
type datachan chan interface{}

// --- json_outstream ---
type json_outstream struct {
	kernel.Algorithm
	w *os.File
	item_names []string
	out datachan
	ctl chan bool
}

func (self *json_outstream) Initialize() kernel.StatusCode {
	self.MsgDebug("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	self.item_names = self.GetProperty("Items").([]string)

	fname := self.GetProperty("Output").(string)
	self.MsgInfo("output file: [%v]\n", fname)
	self.MsgInfo("items: %v\n", self.item_names)

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := uint32(0666)
	w,err := os.Open(fname, flag, perm)
	if err != nil {
		self.MsgError("problem while opening file [%v]: %v\n", fname, err)
		return kernel.StatusCode(1)
	}
	self.w = w
	self.out, self.ctl = data_sink(json.NewEncoder(self.w))
	return kernel.StatusCode(0)
}

func (self *json_outstream) Execute(ctx kernel.IEvtCtx) kernel.StatusCode {
	self.MsgDebug("== execute ==\n")
	store := self.EvtStore(ctx)
	if store == nil {
		self.MsgError("could not retrieve evt-store\n")
	}

	hdr_offset := 1

	val := make([]interface{}, len(self.item_names)+hdr_offset)
	val[0] = ctx.Idx()

	for i,k := range self.item_names {
		val[i+hdr_offset] = store.Get(k)
	}

	self.out <- val
	/*
	err := self.enc.Encode(val)
	if err != nil {
		self.MsgError("error while writing store content: %v\n", err)
		return kernel.StatusCode(1)
	}
	 */
	return kernel.StatusCode(0)
}

func (self *json_outstream) Finalize() kernel.StatusCode {
	self.MsgDebug("== finalize ==\n")
	// close out our data channels
	self.ctl <- true
	close(self.ctl)
	close(self.out)

	self.w.Close()

	return kernel.StatusCode(0)
}

// check implementations match interfaces
var _ = kernel.IComponent(&gob_outstream{})
var _ = kernel.IAlgorithm(&gob_outstream{})

var _ = kernel.IComponent(&json_outstream{})
var _ = kernel.IAlgorithm(&json_outstream{})

// --- factory ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "gob_outstream":
		self := &gob_outstream{}
		kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		// properties
		self.DeclareProperty("Output", "foo.gob")
		self.DeclareProperty("Items", g_keys)
		return self

	case "json_outstream":
		self := &json_outstream{}
		kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		// properties
		self.DeclareProperty("Output", "foo.json")
		self.DeclareProperty("Items", g_keys)
		return self

	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
