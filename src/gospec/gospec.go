// Copyright (c) 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector";
//	"fmt";
)


// Context

type Context struct {
	targetPath []int;
	currentPath *vector.IntVector;
	currentSiblingPos int;
}

func newInitialContext() *Context {
	return newExplicitContext([]int{});
}

func newExplicitContext(targetPath []int) *Context {
	c := new(Context);
	c.targetPath = targetPath;
	c.currentPath = vector.NewIntVector(5);
	c.currentSiblingPos = 0;
	return c;
}

func (c *Context) Specify(name string, closure func()) {
	c.enterChildSpec();
	if c.shouldExecuteCurrentChild() {
		c.executeCurrentChild(closure);
	}
	c.exitChildSpec();
}

func (c *Context) enterChildSpec() {
}

func (c *Context) shouldExecuteCurrentChild() bool {
//	fmt.Printf("targetPath: %v\n", c.targetPath);
//	fmt.Printf("currentPath: %v\n", c.currentPath);
//	fmt.Printf("currentSiblingPos: %v\n", c.currentSiblingPos);
	
	if c.currentSiblingPos == 0 {
		return true;
	}
	return false;
}

func (c *Context) executeCurrentChild(closure func()) {
	closure();
}

func (c *Context) exitChildSpec() {
	c.currentSiblingPos++;
}


// Spec Runner

type RootSpecRunner struct {
	closure func(*Context);
}

func (self *RootSpecRunner) runInContext(c *Context) {
	self.closure(c);
}

