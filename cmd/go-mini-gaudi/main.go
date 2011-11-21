package main

import (
	"flag"
	//"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// thread safe output
var fmt = log.New(os.Stdout, "", log.Lmicroseconds)

var g_do_seq *bool = flag.Bool("seq", false, "enable sequential alg seq")
var g_nprocs *int  = flag.Int("nprocs", 0, "number of threads")
var g_evtmax *int  = flag.Int("evtmax", 10, "number of events")

// type datastore struct {
// 	store map[string]chan interface{}
// }

type datastore map[string]interface{}

func new_datastore() datastore {
	o := make(datastore)
	return o
}

type depgraph map[string]chan int

func new_depgraph() depgraph {
	o := make(depgraph)
	return o
}

type alg struct {
	name  string
	sleep int
	depg  depgraph
	deps  []string
	store datastore
}

func (a *alg) Initialize() error {
	return nil
}

func (a *alg) Execute(ctx int) error {
	for _, dep := range a.deps {
		fmt.Printf(":: [%s] waiting for [%s] (evt: %d)...\n", a.name, dep, ctx)
		v := <-a.depg[dep]
		a.depg[dep] <- v
	}

	// simulate work...
	v := <-time.After(int64(a.sleep) * int64(1e9))
	// 

	fmt.Printf(":: [%s] done (evt: %d) (%v)\n", a.name, ctx, v)
	a.depg[a.name] <- 1

	return nil
}

func (a *alg) Finalize() error {
	return nil
}

func init() {
	flag.Parse()
}

func main() {

	runtime.GOMAXPROCS(*g_nprocs)

	fmt.Printf("== hello ==\n")

	names := []string{
		"alg0",
		"alg1", "alg2", "alg3",
		"alg4", "alg5", "alg6",
		"alg7", "alg8",
		"alg9",
	}
	algs := []alg{}

	depg := new_depgraph()
	store := new_datastore()

	for i, n := range names {
		algs = append(algs, alg{n, i, depg, nil, store})
		depg[n] = make(chan int, 1)
	}
	// reduce sequential part:
	algs[7].sleep = 2
	algs[8].sleep = 2
	algs[9].sleep = 1

	// prepare store layout
	store["017"] = nil
	store["027"] = nil
	store["037"] = nil

	algs[1].deps = []string{"alg0"}
	algs[2].deps = []string{"alg0"}
	algs[3].deps = []string{"alg0"}

	algs[4].deps = []string{"alg0"}
	algs[5].deps = []string{"alg0"}
	algs[6].deps = []string{"alg0"}

	algs[7].deps = []string{"alg1", "alg2", "alg3"}
	algs[8].deps = []string{"alg4", "alg5", "alg6"}

	algs[9].deps = []string{"alg7", "alg8"}

	for _, a := range algs {
		fmt.Printf("--> [init] alg[%s]...\n", a.name)
		if err := a.Initialize(); err != nil {
			panic(err)
		}
	}

	if *g_do_seq {
		for ievt := 0; ievt < *g_evtmax; ievt++ {
			for i, _ := range algs {
				func(ievt, ialg int) {
					a := algs[ialg]
					fmt.Printf("--> [exec-%d] alg[%s]...\n", ievt, a.name)
					if err := a.Execute(ievt); err != nil {
						panic(err)
					}
				}(ievt, i)
			}

			// re-init depg
			for k, _ := range depg {
				depg[k] = make(chan int, 1)
			}
			// re-init store
			for k, _ := range store {
				store[k] = nil
			}
		}
		
	} else {
		for ievt := 0; ievt < *g_evtmax; ievt++ {
			var seq sync.WaitGroup
			seq.Add(len(algs))
			for i, _ := range algs {
				go func(ievt, ialg int) {
					a := algs[ialg]
					fmt.Printf("--> [exec-%d] alg[%s]...\n", ievt, a.name)
					if err := a.Execute(ievt); err != nil {
						panic(err)
					}
					seq.Done()
				}(ievt, i)
			}
			seq.Wait()

			// re-init depg
			for k, _ := range depg {
				depg[k] = make(chan int, 1)
			}
			// re-init store
			for k, _ := range store {
				store[k] = nil
			}
		}
	}

	for _, a := range algs {
		fmt.Printf("--> [fini] alg[%s]...\n", a.name)
		if err := a.Finalize(); err != nil {
			panic(err)
		}
	}

	fmt.Printf("== bye.\n")
}
