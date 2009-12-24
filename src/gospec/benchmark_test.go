// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"testing"
)

const (
	SPECS_COUNT  = 1000
	PRINT_REPORT = false
)

func Benchmark__Running_all_specs(b *testing.B) {
	runSpecs()
}

// TODO: optimize second: much slower than the overhead of running specs (in real life the difference might not be as much)
func Benchmark__Compiling_results(b *testing.B) {
	b.StopTimer()
	runner := runSpecs()
	b.StartTimer()

	runner.compileResults()
}

// TODO: optimize first: this is the slowest, probably because of the string concatenation - print directly to a stream
func Benchmark__Building_a_report(b *testing.B) {
	b.StopTimer()
	runner := runSpecs()
	results := runner.compileResults()
	b.StartTimer()

	buildReport(results)
}


func runSpecs() *Runner {
	runner := NewRunner()
	for i := 0; i < SPECS_COUNT; i++ {
		runner.AddSpec(fmt.Sprintf("DummySpecForBenchmarks%v", i), DummySpecForBenchmarks)
	}
	runner.Run()
	return runner
}

func buildReport(results *ResultCollector) {
	total := results.TotalCount()
	pass := results.PassCount()
	fail := results.FailCount()

	printer := newReportPrinter()
	results.Visit(printer)
	s := printer.String()

	if PRINT_REPORT {
		fmt.Print(s)
		fmt.Printf("Total %v, Pass %v, Fail %v\n", total, pass, fail)
	}
}

func DummySpecForBenchmarks(c *Context) {
	c.Specify("Child A", func() {
		c.Specify("Child AA", func() {
		})
		c.Specify("Child AB", func() {
		})
	})
	c.Specify("Child B", func() {
		c.Specify("Child BA", func() {
		})
		c.Specify("Child BB", func() {
		})
		c.Specify("Child BC", func() {
		})
	})
	c.Specify("Child C", func() {
		c.Specify("Child CA", func() {
		})
		c.Specify("Child CB", func() {
		})
	})
	c.Specify("Child D", func() {
		c.Specify("Child DA", func() {
		})
		c.Specify("Child DB", func() {
		})
		c.Specify("Child DC", func() {
		})
	})
}

