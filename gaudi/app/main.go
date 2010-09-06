package main

import "fmt"
import "flag"
import "os"
import "container/vector"

import "reflect"

import "gaudi/kernel"

var (
	bin  = os.Getenv("GOBIN")
	arch = map[string]string{
		"amd64": "6",
		"386":   "8",
		"arm":   "5",
	}[os.Getenv("GOARCH")]
)

type jobOptions struct {
	pkgs *vector.StringVector
	defs *vector.StringVector
	code *vector.Vector
	exec string
}

var jobOptName *string = flag.String("jobo", "jobOptions.igo",
	"path to a jobOptions file")

func handle_icomponent(c kernel.IComponent) {
	fmt.Printf(":: handle_icomponent(%s)...\n", c.CompName())
}

func main() {
	flag.Parse()

	fmt.Print("::: gaudi\n")
	fmt.Printf("::: getting options from [%s]...\n", *jobOptName)

	app := kernel.NewAppMgr()
	iapp := app.(kernel.IComponent)
	fmt.Printf(" -> created [%s/%s]\n", 
		iapp.CompType(), 
		iapp.CompName())

	{
		t := reflect.Typeof(iapp)
		fmt.Printf("type of [%s]\n", t)
		newt := reflect.NewValue(app)
		fmt.Printf("type of    t: [%s] pkg: [%s] name: [%s]\n", 
			reflect.Typeof(t),
			t.PkgPath(),
			t.Name())
		fmt.Printf("type of newt: [%s]\n", reflect.Typeof(newt))
	}
	handle_icomponent(app)
	fmt.Printf("%s\n", app)

	println("::: configure...")
	sc := app.Configure()
	println("::: configure... [", sc, "]")
	println("::: run...")
	sc = app.Run()
	println("::: run... [", sc, "]")
	fmt.Print("::: bye.\n")
}
