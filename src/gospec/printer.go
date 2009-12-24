// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"io"
)


// Printer formats the spec results into a human-readable format.
type Printer struct {
	out io.Writer
}

func newPrinter(out io.Writer) *Printer {
	return &Printer{out}
}

func (this *Printer) VisitSpec(name string, nestingLevel int, errors []string) {
	// TODO: make the print format pluggable. use this simple version only in tests.
	indent := indent(nestingLevel)
	fmt.Fprintf(this.out, "%v- %v", indent, name)
	if len(errors) > 0 {
		fmt.Fprint(this.out, " [FAIL]\n")
	}
	for _, error := range errors {
		fmt.Fprintf(this.out, "%v    %v", indent, error)
	}
	fmt.Fprint(this.out, "\n")
}

func indent(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += "  "
	}
	return s
}

