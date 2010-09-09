app.algs += Alg("tests/pkg1:Alg1", "alg1", OutputLevel=1)
app.algs += Alg("tests/pkg1:Alg2", "alg2", OutputLevel=1)
app.algs += Alg("tests/pkg2:Alg1", "alg_one", OutputLevel=1)

app.svcs += Svc("tests/pkg1:Svc1", name="svc1", OutputLevel=1)
app.svcs += Svc("tests/pkg2:Svc2", "svc2")

app.toolsvc += Tool("tests/pkg1:Tool1", name="tool1")

