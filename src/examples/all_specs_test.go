// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"github.com/orfjackal/gospec/src/gospec"
	"testing"
)

// You will need to list every spec in a TestXxx method like this,
// so that gotest can be used to run the specs. Later GoSpec might
// get its own command line tool similar to gotest, but for now this
// is the way to go. This shouldn't require too much typing, because
// there will be typically only one top-level spec per class/feature.

func TestAllSpecs(t *testing.T) {
	r := gospec.NewRunner()

	// List all specs here
	r.AddSpec(ExecutionModelSpec)
	r.AddSpec(ExpectationSyntaxSpec)
	r.AddSpec(FibSpec)
	r.AddSpec(StackSpec)

	// Run GoSpec and report any errors to gotest's `testing.T` instance
	gospec.MainGoTest(r, t)
}
