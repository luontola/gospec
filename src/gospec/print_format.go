// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"io"
)


type PrintFormat interface {
	PrintPassing(nestingLevel int, name string)
	PrintFailing(nestingLevel int, name string, errors []*Error)
	PrintSummary(passCount int, failCount int)
}


// PrintFormat for production use.
func DefaultPrintFormat(out io.Writer) PrintFormat {
	return &defaultPrintFormat{out}
}

type defaultPrintFormat struct {
	out io.Writer
}

func (this *defaultPrintFormat) PrintPassing(nestingLevel int, name string) {
	if nestingLevel == 0 {
		fmt.Fprintf(this.out, "\n%v\n", name)
	} else {
		indent := indent(nestingLevel)
		fmt.Fprintf(this.out, "%v- %v\n", indent, name)
	}
}

func (this *defaultPrintFormat) PrintFailing(nestingLevel int, name string, errors []*Error) {
	indent := indent(nestingLevel)
	// TODO: use colors (red)
	fmt.Fprintf(this.out, "%v- %v [FAIL]\n\n", indent, name)
	for _, error := range errors {
		fmt.Fprintf(this.out, "*** %v\n    at %v\n\n", error.Message, error.Location)
	}
	fmt.Fprint(this.out, "\n")
}

func (this *defaultPrintFormat) PrintSummary(passCount int, failCount int) {
	totalCount := passCount + failCount
	// TODO: use colors (red if failures, else green)
	fmt.Fprintf(this.out, "\n%v specs, %v failures\n", totalCount, failCount)
}


// PrintFormat for use in only tests. Does not print line numbers, colors or
// other fancy stuff. Makes comparing as a string easier.
func SimplePrintFormat(out io.Writer) PrintFormat {
	return &simplePrintFormat{out}
}

type simplePrintFormat struct {
	out io.Writer
}

func (this *simplePrintFormat) PrintPassing(nestingLevel int, name string) {
	indent := indent(nestingLevel)
	fmt.Fprintf(this.out, "%v- %v\n", indent, name)
}

func (this *simplePrintFormat) PrintFailing(nestingLevel int, name string, errors []*Error) {
	indent := indent(nestingLevel)
	fmt.Fprintf(this.out, "%v- %v [FAIL]\n", indent, name)
	for _, error := range errors {
		fmt.Fprintf(this.out, "%v    %v\n", indent, error.Message)
	}
}

func (this *simplePrintFormat) PrintSummary(passCount int, failCount int) {
	totalCount := passCount + failCount
	fmt.Fprintf(this.out, "\n%v specs, %v failures\n", totalCount, failCount)
}

func indent(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += "  "
	}
	return s
}

