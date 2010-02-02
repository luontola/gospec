// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
)


// Context controls the execution of the current spec. Child specs can be
// created with the Specify method.
type Context interface {

	// Creates a child spec for the currently executing spec. Specs can be nested
	// unlimitedly. The name should describe what is the behaviour being specified
	// by this spec, and the closure should contain the same specification written
	// as code.
	Specify(name string, closure func())
	
	// Then method starts an assertion. Example:
	//    c.Then(actual).Should.Equal(expected);
	Then(actual interface{}) *MatcherBuilder
}


type taskContext struct {
	targetPath     path
	currentSpec    *specRun
	executedSpecs  *list.List
	postponedSpecs *list.List
}

func newInitialContext() *taskContext {
	return newExplicitContext(rootPath())
}

func newExplicitContext(targetPath path) *taskContext {
	c := new(taskContext)
	c.targetPath = targetPath
	c.currentSpec = nil
	c.executedSpecs = list.New()
	c.postponedSpecs = list.New()
	return c
}

func (c *taskContext) Specify(name string, closure func()) {
	c.enterSpec(name, closure)
	c.processCurrentSpec()
	c.exitSpec()
}

func (c *taskContext) enterSpec(name string, closure func()) {
	spec := newSpecRun(name, closure, c.currentSpec, c.targetPath)
	c.currentSpec = spec
}

func (c *taskContext) processCurrentSpec() {
	spec := c.currentSpec
	switch {
	case c.shouldExecute(spec):
		c.execute(spec)
	case c.shouldPostpone(spec):
		c.postpone(spec)
	}
}

func (c *taskContext) exitSpec() {
	c.currentSpec = c.currentSpec.parent
}

func (c *taskContext) shouldExecute(spec *specRun) bool {
	if spec.parent != nil && spec.parent.hasFatalErrors {
		return false
	}
	return spec.isOnTargetPath() || (spec.isUnseen() && spec.isFirstChild())
}

func (c *taskContext) shouldPostpone(spec *specRun) bool {
	return spec.isUnseen() && !spec.isFirstChild()
}

func (c *taskContext) execute(spec *specRun) {
	c.executedSpecs.PushBack(spec)
	spec.execute()
}

func (c *taskContext) postpone(spec *specRun) {
	c.postponedSpecs.PushBack(spec)
}

func (c *taskContext) Then(actual interface{}) *MatcherBuilder {
	return newMatcherBuilder(actual, callerLocation(), c.currentSpec)
}

