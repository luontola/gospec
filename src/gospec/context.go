// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list";
)


// Context controls the execution of the current spec. Child specs can be
// created with the Specify method.
type Context struct {
	targetPath path;
	currentSpec *specification;
	executedSpecs *list.List;
	postponedSpecs *list.List;
	done chan bool;
}

func newInitialContext() *Context {
	return newExplicitContext(rootPath())
}

func newExplicitContext(targetPath path) *Context {
	c := new(Context);
	c.targetPath = targetPath;
	c.currentSpec = nil;
	c.executedSpecs = list.New();
	c.postponedSpecs = list.New();
	c.done = make(chan bool);
	return c
}

// Creates a child spec for the currently executing spec. Specs can be nested
// unlimitedly. The name should describe what is the behaviour being specified
// by this spec, and the closure should contain the same specification written
// as code.
func (c *Context) Specify(name string, closure func()) {
	c.enterSpec(name, closure);
	c.processCurrentSpec();
	c.exitSpec();
}

func (c *Context) enterSpec(name string, closure func()) {
	spec := newSpecification(name, closure, c.currentSpec);
	c.currentSpec = spec;
}

func (c *Context) processCurrentSpec() {
	spec := c.currentSpec;
	switch {
	case c.shouldExecute(spec):
		c.execute(spec)
	case c.shouldPostpone(spec):
		c.postpone(spec)
	}
}

func (c *Context) exitSpec() {
	c.currentSpec = c.currentSpec.parent;
}

func (c *Context) shouldExecute(spec *specification) bool {
	return spec.isOnTargetPath(c) || (spec.isUnseen(c) && spec.isFirstChild())
}

func (c *Context) shouldPostpone(spec *specification) bool {
	return spec.isUnseen(c) && !spec.isFirstChild()
}

func (c *Context) execute(spec *specification) {
	c.executedSpecs.PushBack(spec);
	spec.execute();
}

func (c *Context) postpone(spec *specification) {
	c.postponedSpecs.PushBack(spec);
}

