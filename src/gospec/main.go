// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"flag"
	"os"
	"testing"
)


var (
    printAll = flag.Bool("print-all", false, "print also passing specs and not only failing (GoSpec)")
)

// Executes the specs which have been added to the Runner
// and prints the results to stdout.
func MainGoTest(runner *Runner, t *testing.T) {
	// Assume that this method will then be executed by gotest and
	// flag.Parse() has already been called in testing.Main() so 
	// we don't need to call it here.
	
	printer := NewPrinter(DefaultPrintFormat(os.Stdout))
	if *printAll {
		printer.ShowAll()
	} else {
		printer.ShowOnlyFailing()
	}
	printer.ShowSummary()
	
	runner.Run()
	results := runner.compileResults()
	results.Visit(printer)
	
	// TODO: FailCount() is not optimized - it visits all specs again
	if results.FailCount() > 0 {
		t.Fail()
	}
}

