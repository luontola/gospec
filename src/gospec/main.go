// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
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
// and prints the results to stdout. Exits the process after
// it is finished - with zero or non-zero exit value,
// depending on whether any specs failed.
func Main(runner *Runner) {
	flag.Parse()
	results := runAndPrint(runner)
	if results.FailCount() > 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Executes the specs which have been added to the Runner
// and prints the results to stdout. Fails the surrounding
// test if any of the specs fails.
func MainGoTest(runner *Runner, t *testing.T) {
	// Assume that this method will then be executed by gotest and
	// flag.Parse() has already been called in testing.Main() so
	// we don't need to call it here.

	results := runAndPrint(runner)
	if results.FailCount() > 0 {
		t.Fail()
	}
}

func runAndPrint(runner *Runner) *ResultCollector {
	printer := NewPrinter(DefaultPrintFormat(os.Stdout))
	if *printAll {
		printer.ShowAll()
	} else {
		printer.ShowOnlyFailing()
	}
	printer.ShowSummary()

	runner.Run()
	results := runner.Results()
	results.Visit(printer)
	return results
}
