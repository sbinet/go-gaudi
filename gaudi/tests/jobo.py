app.algs += Alg("gaudi/tests/pkg1:Alg1", "alg1", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg1:Alg2", "alg2", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg2:Alg1", "alg_one", OutputLevel=1)

app.svcs += Svc("gaudi/tests/pkg1:Svc1", name="svc1", OutputLevel=1)
app.svcs += Svc("gaudi/tests/pkg2:Svc2", "svc2")

app.toolsvc += Tool("gaudi/tests/pkg1:Tool1", name="tool1")

