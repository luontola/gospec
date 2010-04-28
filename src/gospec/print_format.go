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
		fmt.Fprintf(this.out, "%v- %v\n", indent(nestingLevel), name)
	}
}

func (this *defaultPrintFormat) PrintFailing(nestingLevel int, name string, errors []*Error) {
	// TODO: use colors (red)
	fmt.Fprintf(this.out, "%v- %v [FAIL]\n\n", indent(nestingLevel), name)
	for _, error := range errors {
		this.printError(error)
	}
	fmt.Fprint(this.out, "\n")
}

func (this *defaultPrintFormat) printError(error *Error) {
	// Go's stack trace format can be seen in
	// traceback() at src/pkg/runtime/amd64/traceback.c
	// but we don't have to use exactly the same format.
	fmt.Fprint(this.out, formatErrorMessage(error))
	for _, loc := range error.StackTrace {
		fmt.Fprintf(this.out, "    %v()  at  %v:%v\n", loc.Name(), loc.File(), loc.Line())
	}
	fmt.Fprintf(this.out, "\n")
}

func formatErrorMessage(e *Error) string {
	s := ""
	switch e.Type {
	case ExpectFailed:
		s += fmt.Sprintf("*** Expected: %v\n", e.Message)
		s += fmt.Sprintf("         got: “%v”\n", e.Actual)
	case AssumeFailed:
		s += fmt.Sprintf("*** Assumed: %v\n", e.Message)
		s += fmt.Sprintf("        got: “%v”\n", e.Actual)
	case OtherError:
		s += fmt.Sprintf("*** %v\n", e.Message)
	}
	return s
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
	fmt.Fprintf(this.out, "%v- %v\n", indent(nestingLevel), name)
}

func (this *simplePrintFormat) PrintFailing(nestingLevel int, name string, errors []*Error) {
	fmt.Fprintf(this.out, "%v- %v [FAIL]\n", indent(nestingLevel), name)
	for _, error := range errors {
		this.printError(error)
	}
}

func (this *simplePrintFormat) printError(error *Error) {
	fmt.Fprintf(this.out, formatErrorMessage(error))
	for _, loc := range error.StackTrace {
		fmt.Fprintf(this.out, "    at %v\n", loc.FileName())
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
