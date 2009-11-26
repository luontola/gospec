// Copyright (c) 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
)


// Common test utils

var testSpy = "";

func resetTestSpy() {
	testSpy = "";
}

func assertTestSpyHas(expected string, t *testing.T) {
	assertEquals(expected, testSpy, t);
}

func assertEquals(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Error("Expected '" + expected + "' but was '" + actual + "'");
	}
}

