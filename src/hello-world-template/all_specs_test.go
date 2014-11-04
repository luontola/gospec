package hello

import (
	"github.com/orfjackal/gospec/src/gospec"
	"testing"
)

func TestAllSpecs(t *testing.T) {
	r := gospec.NewParallelRunner()
	r.AddSpec(HelloSpec)
	gospec.MainGoTest(r, t)
}
