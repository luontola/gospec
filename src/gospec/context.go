// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
)

// Context controls the execution of the current spec. Child specs can be
// created with the Specify method.
type Context interface {

	// Creates a child spec for the currently executing spec. Specs can be
	// nested unlimitedly. The name should describe what is the behaviour being
	// specified by this spec, and the closure should express the same
	// specification as code.
	Specify(name string, closure func())

	// Makes an expectation. For example:
	//    c.Expect(theAnswer, Equals, 42)
	//    c.Expect(theAnswer, Not(Equals), 666)
	//    c.Expect(thereIsASpoon, IsFalse)
	Expect(actual interface{}, matcher Matcher, expected ...interface{})

	// Makes an assumption. Otherwise the same as an expectation,
	// but on failure will not continue executing the child specs.
	Assume(actual interface{}, matcher Matcher, expected ...interface{})
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

func (c *taskContext) Expect(actual interface{}, matcher Matcher, expected ...interface{}) {
	location := callerLocation()
	logger := expectationLogger{c.currentSpec}
	m := newMatcherAdapter(location, logger, ExpectFailed)
	m.Expect(actual, matcher, expected...)
}

func (c *taskContext) Assume(actual interface{}, matcher Matcher, expected ...interface{}) {
	location := callerLocation()
	logger := assumptionLogger{c.currentSpec}
	m := newMatcherAdapter(location, logger, AssumeFailed)
	m.Expect(actual, matcher, expected...)
}

type expectationLogger struct {
	log ratedErrorLogger
}

func (this expectationLogger) AddError(e *Error) {
	this.log.AddError(e)
}

type assumptionLogger struct {
	log ratedErrorLogger
}

func (this assumptionLogger) AddError(e *Error) {
	this.log.AddFatalError(e)
}
