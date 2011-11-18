# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ${GOROOT}/src/Make.inc

all: install

DIRS=\
	pkg/gaudi/kernel\
	pkg/gaudi/kernel/datastore\
	pkg/gaudi/kernel/evtproc\
	pkg/gaudi/kernel/outstream\
	cmd/go-gaudi\
	pkg/gaudi/tests/pkg1\
	pkg/gaudi/tests/pkg2\


clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
bench.dirs: $(addsuffix .bench, $(BENCH))

%.clean:
	+cd $* && $(QUOTED_GOBIN)/gomake clean

%.install:
	+cd $* && $(QUOTED_GOBIN)/gomake install

clean: clean.dirs

install: install.dirs

#-include ${GOROOT}/src/Make.deps
