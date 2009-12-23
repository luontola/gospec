// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
)


// ReportPrinter formats the spec results into a human-readable format.
type ReportPrinter struct {
	report      string
	indentLevel int
}

func newReportPrinter() *ReportPrinter {
	return &ReportPrinter{"", 0}
}

func (this *ReportPrinter) Visit(results *ResultCollector) {
	for rootSpec := range results.Roots() {
		this.visitSpec(rootSpec)
	}
}

func (this *ReportPrinter) visitSpec(spec *specResult) {
	// TODO: make the print format pluggable. use this simple version only in tests.
	errors := ""
	if spec.IsFailed() {
		errors += " [FAIL]\n"
	}
	for error := range spec.errors.Iter() {
		errors += fmt.Sprintf("%v    %v", this.indent(), error)
	}

	this.report += fmt.Sprintf("%v- %v%v\n", this.indent(), spec.name, errors)

	for child := range spec.Children() {
		this.indentLevel++
		this.visitSpec(child)
		this.indentLevel--
	}
}

func (this *ReportPrinter) indent() string {
	s := ""
	for i := 0; i < this.indentLevel; i++ {
		s += "  "
	}
	return s
}

func (this *ReportPrinter) String() string {
	return this.report
}

