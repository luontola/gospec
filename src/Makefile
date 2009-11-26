include $(GOROOT)/src/Make.$(GOARCH)

#Run 'make deps' if you add or change new folders.
DIRS=\
	gospec\

TEST=\
	gospec\

all.dirs: $(addsuffix .all, $(DIRS))
clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
docs.dirs: $(addsuffix .docs, $(DIRS))

%.all:
	+cd $* && make

%.clean:
	+cd $* && make clean

%.install:
	+cd $* && make install

%.test:
	+cd $* && make test

%.docs:
	+godoc -pkgroot="." -html -v $* > docs/$*.html


clean:	clean.dirs

install: install.dirs

test:	test.dirs

docs:	docs.dirs

deps:
	./deps.bash

include Make.deps

