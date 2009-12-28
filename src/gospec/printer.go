// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"io"
)


type printMode int

const (
	ALL printMode = iota
	ONLY_FAILING
)


// Printer formats the spec results into a human-readable format.
type Printer struct {
	out         io.Writer
	show        printMode
	showSummary bool
	notPrinted  []string
}

func newPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,
		show: ALL,
		showSummary: true,
		notPrinted: []string{},
	}
}

func (this *Printer) ShowAll() {
	this.show = ALL
}

func (this *Printer) ShowOnlyFailing() {
	this.show = ONLY_FAILING
}

func (this *Printer) HideSummary() {
	this.showSummary = false
}

func (this *Printer) ShowSummary() {
	this.showSummary = true
}

func (this *Printer) VisitSpec(nestingLevel int, name string, errors []*Error) {
	isPassing := len(errors) == 0
	isFailing := !isPassing
	
	if isPassing {
		if this.show == ALL {
			this.printPassing(nestingLevel, name)
		} else {
			this.saveNotPrinted(nestingLevel, name)
		}
	}
	if isFailing {
		this.printNotPrintedParents(nestingLevel)
		this.printFailing(nestingLevel, name, errors)
	}
}

func (this *Printer) VisitEnd(passCount int, failCount int) {
	if this.showSummary {
		this.printSummary(passCount, failCount)
	}
}


func (this *Printer) saveNotPrinted(nestingLevel int, name string) {
	if nestingLevel >= len(this.notPrinted) {
		resizeArray(&this.notPrinted, nestingLevel+1)
	}
	this.notPrinted[nestingLevel] = name
}

func (this *Printer) printNotPrintedParents(nestingLevel int) {
	for i, name := range this.notPrinted {
		if i < nestingLevel && name != "" {
			this.printPassing(i, name)
		}
		this.notPrinted[i] = ""
	}
}

func resizeArray(arr *[]string, newLength int) {
	old := *arr
	*arr = make([]string, newLength)
	copy(*arr, old)
}


// TODO: make the print format pluggable. use this simple version only in tests.

func (this *Printer) printPassing(nestingLevel int, name string) {
	indent := indent(nestingLevel)
	fmt.Fprintf(this.out, "%v- %v\n", indent, name)
}

func (this *Printer) printFailing(nestingLevel int, name string, errors []*Error) {
	indent := indent(nestingLevel)
	fmt.Fprintf(this.out, "%v- %v [FAIL]\n", indent, name)
	for _, error := range errors {
		// TODO: print file name and line number
		// example:
		// foo.go:23  Expected X but was Y
		fmt.Fprintf(this.out, "%v    %v\n", indent, error.Message)
	}
}

func (this *Printer) printSummary(passCount int, failCount int) {
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

