package main

import "gaudi/kernel"
import "fmt"


func handle_icomponent(c kernel.IComponent) {
	fmt.Printf(":: handle_icomponent(%s)...\n", c.CompName())
}

func main() {
	fmt.Print("::: gaudi\n")
	app := kernel.NewAppMgr()

	fmt.Printf(" -> created [%s/%s]\n", 
		app.(kernel.IComponent).CompType(), 
		app.(kernel.IComponent).CompName())

	handle_icomponent(app)
	fmt.Printf("%s\n", app)

	println("::: configure...")
	sc := app.Configure()
	println("::: configure... [", sc, "]")
	println("::: run...")
	sc = app.Run()
	println("::: run... [", sc, "]")

	// println("::: testing event server...")
	// e := app.evtproc.(*evtProc)
	// e.test_0()
	// println("::: testing event server... [done]")
	fmt.Print("::: bye.\n")
}
