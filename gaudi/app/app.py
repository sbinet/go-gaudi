#!/usr/bin/env python

"""\
Driver for the ng-gaudi 'framework'.

This will read a python jobo file (or a set of those) and create a go-binary from the assembled set of go-packages and ng-components.
"""

### imports ------------------------------------------------------------------
import sys
import os

### functions ----------------------------------------------------------------
def _make_configurable(pkg_name, name, **kwds):
    """create a configurable for a component from package `pkg_name` with
    name `name`. the keyword arguments are its properties
    """
    kwds['name']=name
    return Configurable(pkg_name=pkg_name, **kwds)
Alg = _make_configurable
Svc = _make_configurable
Tool= _make_configurable

### classes ------------------------------------------------------------------

class Configurable(object):
    def __init__(self, pkg_name, **kwds):
        self.pkg_name = pkg_name
        self.name = kwds['name']
        self.props = dict(kwds)
        del self.props['name']
        
class AppMgr(object):
    def __init__(self):
        self.algs = []
        self.svcs = []
        self.toolsvc = []
        return

    def configure(self, jobopts):
        if not isinstance(jobopts, (list,tuple)):
            raise TypeError("jobopts should be a sequence")
        self.jobopts = jobopts[:]
        for i,jobo in enumerate(self.jobopts):
            jobo = os.path.expandvars(os.path.expanduser(jobo))
            self.jobopts[i] = jobo
            if not os.path.exists(jobo):
                raise RuntimeError("no such file [%s]" % jobo)
            
    def run(self):
        return 0
    
app = AppMgr()

if __name__ == "__main__":
    jobopts = []
    if len(sys.argv)>1:
        for arg in sys.argv[1:]:
            if arg.endswith(".py"):
                jobopts += [arg]
    if len(jobopts) == 0:
        jobopts = ["jobopt.py"]
    app.configure(jobopts)
    sys.exit(app.run())
