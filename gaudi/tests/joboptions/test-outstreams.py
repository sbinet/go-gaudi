app.props.EvtMax = 50
app.props.OutputLevel = 1
#app.props.NbrProcs = 1

nworkers = 50
#nworkers = 1

app.svcs += Svc("gaudi/kernel/evtproc:evtproc",
                "evt-proc",
                OutputLevel=Lvl.INFO,
                NbrWorkers=nworkers)

app.svcs += Svc("gaudi/kernel/datastore:datastoresvc", "evt-store")
app.svcs += Svc("gaudi/kernel/datastore:datastoresvc", "det-store")

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_1",
                OutputLevel=Lvl.ERROR,
                Val=0.)
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-1",
                ExpectedValue=1)

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_2",
                OutputLevel=Lvl.ERROR,
                Val=3.)
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-2",
                ExpectedValue=2)

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_3")
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-3",
                ExpectedValue=3)

app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper",
                NbrJets = "njets",
                ExpectedValue=3,
                OutputLevel=Lvl.ERROR)

app.algs += Alg("gaudi/kernel/outstream:gob_outstream",
                "gobwriter",
                Items=["cnt",],
                Output="/tmp/foo.gob")

app.algs += Alg("gaudi/kernel/outstream:json_outstream",
                "json-writer-1",
                Output="/tmp/foo.1.json")

app.algs += Alg("gaudi/kernel/outstream:json_outstream",
                "json-writer-2",
                Output="/tmp/foo.2.json")

app.toolsvc += Tool("gaudi/tests/pkg1:tool1", name="tool1")

