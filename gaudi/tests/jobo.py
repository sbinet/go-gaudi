app.algs += Alg("gaudi/tests/pkg1:alg1", "alg1", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg1:alg2", "alg2", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg2:alg1", "alg_one", OutputLevel=1)

app.svcs += Svc("gaudi/tests/pkg1:svc1", name="svc1", OutputLevel=1)
app.svcs += Svc("gaudi/tests/pkg2:svc2", "svc2")

app.svcs += Svc("gaudi/datastore:datastoresvc", "evt-store")
app.svcs += Svc("gaudi/datastore:datastoresvc", "det-store")

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_1",
                OutputLevel=1,
                Val=0.)
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-1",
                ExpectedValue=1)

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_2",
                OutputLevel=1,
                Val=3.)
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-2",
                ExpectedValue=2)

app.algs += Alg("gaudi/tests/pkg2:alg_adder", "adder_3")
app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper-3",
                ExpectedValue=3)

app.algs += Alg("gaudi/tests/pkg2:alg_dumper", "dumper",
                NbrJets = "njets",
                ExpectedValue=3,
                OutputLevel=1)

for i in xrange(1000):
    app.algs += Alg("gaudi/tests/pkg2:alg_adder",
                    "adder--%04i" % i,
                    SimpleCounter="my_counter")
    app.algs += Alg("gaudi/tests/pkg2:alg_dumper",
                    "dumper--%04i" % i,
                    SimpleCounter="my_counter",
                    ExpectedValue=i+1)
    
app.toolsvc += Tool("gaudi/tests/pkg1:tool1", name="tool1")

