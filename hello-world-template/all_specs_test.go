package hello

import (
	"gospec"
	"testing"
)


func TestAllSpecs(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(HelloSpec)
	gospec.MainGoTest(r, t)
}
