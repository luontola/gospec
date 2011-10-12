// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import ()

type printMode int

const (
	ALL printMode = iota
	ONLY_FAILING
)

// Printer formats the spec results into a human-readable format.
type Printer struct {
	format      PrintFormat
	show        printMode
	showSummary bool
	notPrinted  []string
}

func NewPrinter(format PrintFormat) *Printer {
	return &Printer{
		format:      format,
		show:        ALL,
		showSummary: true,
		notPrinted:  []string{},
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
			this.format.PrintPassing(nestingLevel, name)
		} else {
			this.saveNotPrinted(nestingLevel, name)
		}
	}
	if isFailing {
		this.printNotPrintedParents(nestingLevel)
		this.format.PrintFailing(nestingLevel, name, errors)
	}
}

func (this *Printer) VisitEnd(passCount int, failCount int) {
	if this.showSummary {
		this.format.PrintSummary(passCount, failCount)
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
			this.format.PrintPassing(i, name)
		}
		this.notPrinted[i] = ""
	}
}

func resizeArray(arr *[]string, newLength int) {
	old := *arr
	*arr = make([]string, newLength)
	copy(*arr, old)
}
