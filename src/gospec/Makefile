include $(GOROOT)/src/Make.inc

TARG=gospec
GOFILES=\
	context.go\
	error.go\
	funcname.go\
	location.go\
	main.go\
	matchers.go\
	printer.go\
	print_format.go\
	recover.go\
	results.go\
	runner.go\
	specification.go\

# Must disable inlining, or tests will fail because of stack traces
# not containing the innermost method call.
GCFLAGS+=-l

include $(GOROOT)/src/Make.pkg

