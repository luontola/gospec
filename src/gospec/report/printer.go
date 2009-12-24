// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package report

import (
	"fmt"
)


// ReportPrinter formats the spec results into a human-readable format.
type ReportPrinter struct {
	report      string
}

func newReportPrinter() *ReportPrinter {
	return &ReportPrinter{""}
}

func (this *ReportPrinter) VisitSpec(name string, nestingLevel int, errors []string) {
	// TODO: make the print format pluggable. use this simple version only in tests.
	s := ""
	if len(errors) > 0 {
		s += " [FAIL]\n"
	}
	for _, error := range errors {
		s += fmt.Sprintf("%v    %v", indent(nestingLevel), error)
	}
	
	this.report += fmt.Sprintf("%v- %v%v\n", indent(nestingLevel), name, s)
}

func indent(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += "  "
	}
	return s
}

func (this *ReportPrinter) String() string {
	return this.report
}

