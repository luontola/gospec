// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

const (
	SPECS_COUNT  = 10000
	PRINT_REPORT = false
)

// 2009-12-25: Compiling results and building reports takes a long time
// when using GOMAXPROCS=4 under a virtual machine with only one CPU.
// TODO: Run the benchmarks on native hardware, try different values of GOMAXPROCS, use 6prof.

func Benchmark__Running_all_specs(b *testing.B) {
	runSpecs()
}

func Benchmark__Compiling_results(b *testing.B) {
	b.StopTimer()
	runner := runSpecs()
	b.StartTimer()

	runner.compileResults()
}

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
	
	var report io.Writer
	if PRINT_REPORT {
		report = new(bytes.Buffer)
	} else {
		report = new(NullWriter)
	}
	
	results.Visit(newPrinter(report))

	if PRINT_REPORT {
		buf := report.(*bytes.Buffer)
		buf.WriteTo(os.Stdout)
		fmt.Printf("Total %v, Pass %v, Fail %v\n", total, pass, fail)
	}
}


type NullWriter struct {}

func (w *NullWriter) Write(p []byte) (n int, err os.Error) {
	return len(p), nil
}


func DummySpecForBenchmarks(c *Context) {
	c.Specify("Child A", func() {
		c.Specify("Child AA", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child AB", func() {
			c.Then(1).Should.Equal(1)
		})
	})
	c.Specify("Child B", func() {
		c.Specify("Child BA", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child BB", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child BC", func() {
			c.Then(1).Should.Equal(1)
		})
	})
	c.Specify("Child C", func() {
		c.Specify("Child CA", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child CB", func() {
			c.Then(1).Should.Equal(1)
		})
	})
	c.Specify("Child D", func() {
		c.Specify("Child DA", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child DB", func() {
			c.Then(1).Should.Equal(1)
		})
		c.Specify("Child DC", func() {
			c.Then(1).Should.Equal(1)
		})
	})
}

