app.algs += Alg("gaudi/tests/pkg1:alg1", "alg1", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg1:alg2", "alg2", OutputLevel=1)
app.algs += Alg("gaudi/tests/pkg2:alg1", "alg_one", OutputLevel=1)

app.svcs += Svc("gaudi/tests/pkg1:svc1", name="svc1", OutputLevel=1)
app.svcs += Svc("gaudi/tests/pkg2:svc2", "svc2")

app.toolsvc += Tool("gaudi/tests/pkg1:tool1", name="tool1")

