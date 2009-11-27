// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector";
	"fmt";
)


// Context coordinates the spec execution

type Context struct {
	targetPath []int;
	currentSpec *specification;
	specs *vector.Vector;
}

func newInitialContext() *Context {
	return newExplicitContext([]int{})
}

func newExplicitContext(targetPath []int) *Context {
	c := new(Context);
	c.targetPath = targetPath;
	c.currentSpec = nil;
	c.specs = vector.New(0);
	return c
}

func (c *Context) Specify(name string, closure func()) {
	c.enterSpec(name, closure);
	if c.shouldExecuteCurrentSpec() {
		c.executeCurrentSpec();
	}
	c.exitSpec();
}

func (c *Context) enterSpec(name string, closure func()) {
	spec := newSpecification(name, closure, c.currentSpec);
	c.specs.Push(spec);
	c.currentSpec = spec;
}

func (c *Context) shouldExecuteCurrentSpec() bool {
//	fmt.Println();
//	fmt.Println("targetPath:", c.targetPath);
//	fmt.Println("currentSpec:", c.currentSpec);
	
	isBelowTargetPath := currentIsBelowTargetPath(c.currentSpec.path, c.targetPath);
	isUnseen := len(c.currentSpec.path) > len(c.targetPath);
	isFirstChild := c.currentSpec.lastPathIndex() == 0;
	
	return isBelowTargetPath || (isUnseen && isFirstChild)
}

func currentIsBelowTargetPath(current []int, target []int) bool {
	if len(current) > len(target) {
		return false
	}
	for i := 0; i < len(current); i++ {
		if current[i] != target[i] {
			return false
		}
	}
	return true
}

func (c *Context) executeCurrentSpec() {
	c.currentSpec.execute();
}

func (c *Context) exitSpec() {
	c.currentSpec = c.currentSpec.parent;
}


// Represents a spec in a tree of specs

type specification struct {
	name string;
	closure func();
	parent *specification;
	numberOfChildren int;
	path []int;
}

func newSpecification(name string, closure func(), parent *specification) *specification {
	path := []int{};
	if parent != nil {
		path = pathForChildOf(parent.path, parent.numberOfChildren);
		parent.numberOfChildren++;
	}
	return &specification{name, closure, parent, 0, path}
}

func pathForChildOf(parentPath []int, childIndex int) []int {
	childPath := make([]int, len(parentPath) + 1);
	for i, v := range parentPath {
		childPath[i] = v
	}
	childPath[len(parentPath)] = childIndex;
	return childPath
}

func (s *specification) execute() {
	s.closure();
}

func (s *specification) lastPathIndex() int {
	if len(s.path) == 0 {
		return 0	// root specification
	}
	return s.path[len(s.path) - 1]
}

func (s *specification) String() string {
	return fmt.Sprintf("specification{%v @ %v}", s.name, s.path)
}


// Starts the spec execution

type RootSpecRunner struct {
	specName string;
	specClosure func(*Context);
}

func NewRootSpecRunner(specName string, specClosure func(*Context)) *RootSpecRunner {
	return &RootSpecRunner{specName, specClosure};
}

func (r *RootSpecRunner) runInContext(c *Context) {
	c.Specify(r.specName, func() { r.specClosure(c) });
	
	//fmt.Println(c.specs);
}

