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
	keys := []string{
		"njets",
		"ptjets",
		"cnt",
	}
	allgood := true
	for _,k := range keys {
		err = self.enc.Encode(store.Get("njets"))
		if err != nil {
			self.MsgError("error while writing store content at [%v]: %v\n", 
				k,err)
			allgood = false
		}
	}
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

// --- json_outstream ---
type json_outstream struct {
	kernel.Algorithm
	w *os.File
	enc *json.Encoder
	item_names []string
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

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := uint32(0666)
	w,err := os.Open(fname, flag, perm)
	if err != nil {
		self.MsgError("problem while opening file [%v]: %v\n", fname, err)
		return kernel.StatusCode(1)
	}
	self.w = w
	self.enc = json.NewEncoder(self.w)

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

	err := self.enc.Encode(val)
	if err != nil {
		self.MsgError("error while writing store content: %v\n", err)
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *json_outstream) Finalize() kernel.StatusCode {
	self.MsgDebug("== finalize ==\n")
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
