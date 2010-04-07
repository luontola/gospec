// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
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
	ROOT_SPEC_COUNT = 1000
	PRINT_REPORT    = false
)

// TODO: Run the benchmarks on native hardware, try different values of GOMAXPROCS, use 6prof.

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runSpecs()
	}
}

func BenchmarkResults(b *testing.B) {
	b.StopTimer()
	runner := runSpecs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		runner.Results()
	}
}

func BenchmarkReport(b *testing.B) {
	b.StopTimer()
	runner := runSpecs()
	results := runner.Results()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		buildReport(results)
	}
}


func runSpecs() *Runner {
	runner := NewRunner()
	for i := 0; i < ROOT_SPEC_COUNT; i++ {
		runner.AddNamedSpec(fmt.Sprintf("DummySpecForBenchmarks%v", i), DummySpecForBenchmarks)
	}
	runner.Run()
	return runner
}

func buildReport(results *ResultCollector) {
	var report io.Writer
	if PRINT_REPORT {
		report = new(bytes.Buffer)
	} else {
		report = new(NullWriter)
	}

	results.Visit(NewPrinter(SimplePrintFormat(report)))

	if PRINT_REPORT {
		buf := report.(*bytes.Buffer)
		buf.WriteTo(os.Stdout)
	}
}


type NullWriter struct{}

func (w *NullWriter) Write(p []byte) (n int, err os.Error) {
	return len(p), nil
}


func DummySpecForBenchmarks(c Context) {
	// Some work, to create a more realistic workload and
	// to put the framework's overhead into proportion.
	//for i := 0; i < 1000000; i++ {}

	// 15 spec declarations, executed in 10 runs
	// (each run is 3 levels deep, so in total 30 spec runs)
	c.Specify("Child A", func() {
		c.Specify("Child AA", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child AB", func() {
			c.Expect(1, Equals, 1)
		})
	})
	c.Specify("Child B", func() {
		c.Specify("Child BA", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child BB", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child BC", func() {
			c.Expect(1, Equals, 1)
		})
	})
	c.Specify("Child C", func() {
		c.Specify("Child CA", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child CB", func() {
			c.Expect(1, Equals, 1)
		})
	})
	c.Specify("Child D", func() {
		c.Specify("Child DA", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child DB", func() {
			c.Expect(1, Equals, 1)
		})
		c.Specify("Child DC", func() {
			c.Expect(1, Equals, 1)
		})
	})
}
