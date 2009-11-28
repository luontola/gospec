// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector";
)


// Starts the spec execution

type SpecRunner struct {
	specName string;
	specClosure func(*Context);
	executed *vector.Vector;
	pathsToExecute *vector.Vector;
}

func NewSpecRunner(specName string, specClosure func(*Context)) *SpecRunner {
	executed := vector.New(0);
	pathsToExecute := vector.New(0);
	pathsToExecute.Push(rootPath());
	return &SpecRunner{specName, specClosure, executed, pathsToExecute};
}

func (r *SpecRunner) Run() {
	for r.hasPathsToExecute() {
		r.executeNextPath();
	}
}

func (r *SpecRunner) executeNextPath() {
	targetPath := r.nextPathToExecute();
	context := newExplicitContext(targetPath);
	result := r.runInContext(context);
	r.saveResult(result);
}

func (r *SpecRunner) hasPathsToExecute() bool	{ return r.pathsToExecute.Len() > 0 }
func (r *SpecRunner) nextPathToExecute() path	{ return r.pathsToExecute.Pop().(path) }

func (r *SpecRunner) runInContext(c *Context) *runResult {
	c.Specify(r.specName, func() { r.specClosure(c) });
	return &runResult{
		asSpecArray(c.executedSpecs),
		asSpecArray(c.postponedSpecs)
	}
}

func (r *SpecRunner) saveResult(result *runResult) {
	for _, spec := range result.executedSpecs {
		r.executed.Push(spec);
	}
	for _, spec := range result.postponedSpecs {
		r.pathsToExecute.Push(spec.path);
	}
}

type runResult struct {
	executedSpecs []*specification;
	postponedSpecs []*specification;
}

