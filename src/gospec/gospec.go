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
	targetPath path;
	currentSpec *specification;
	specs *vector.Vector;
}

func newInitialContext() *Context {
	return newExplicitContext(emptyPath())
}

func newExplicitContext(targetPath path) *Context {
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
	
	currentPath := c.currentSpec.path;
	targetPath := c.targetPath;
	
	isOnTargetPath := currentPath.isOn(targetPath);
	isUnseen := currentPath.isBeyond(targetPath);
	isFirstChild := currentPath.lastIndex() == 0;
	
	return isOnTargetPath || (isUnseen && isFirstChild)
}

func (c *Context) executeCurrentSpec() {
	c.currentSpec.execute();
}

func (c *Context) exitSpec() {
	c.currentSpec = c.currentSpec.parent;
}


// Spec paths

type path []int;

func emptyPath() path {
	return []int{};
}

func (parent path) append(index int) path {
	result := make([]int, len(parent) + 1);
	for i, v := range parent {
		result[i] = v
	}
	result[len(parent)] = index;
	return result
}

func (current path) isOn(target path) bool {
	if current.isBeyond(target) {
		return false
	}
	for i := 0; i < len(current); i++ {
		if current[i] != target[i] {
			return false
		}
	}
	return true
}

func (current path) isBeyond(target path) bool {
	return len(current) > len(target)
}

func (path path) lastIndex() int {
	if len(path) == 0 {
		return -1	// empty path, i.e. root specification
	}
	return path[len(path) - 1]
}


// Represents a spec in a tree of specs

type specification struct {
	name string;
	closure func();
	parent *specification;
	numberOfChildren int;
	path path;
}

func newSpecification(name string, closure func(), parent *specification) *specification {
	path := emptyPath();
	if parent != nil {
		currentIndex := parent.numberOfChildren;
		path = parent.path.append(currentIndex);
		parent.numberOfChildren++;
	}
	return &specification{name, closure, parent, 0, path}
}

func (s *specification) execute() {
	s.closure();
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

