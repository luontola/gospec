// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list";
	"container/vector";
	"fmt";
)


// Context coordinates the spec execution

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


// Spec paths

type path []int;

func rootPath() path {
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
		return -1	// root path
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
	path := rootPath();
	if parent != nil {
		currentIndex := parent.numberOfChildren;
		path = parent.path.append(currentIndex);
		parent.numberOfChildren++;
	}
	return &specification{name, closure, parent, 0, path}
}

func (spec *specification) isOnTargetPath(c *Context) bool	{ return spec.path.isOn(c.targetPath) }
func (spec *specification) isUnseen(c *Context) bool		{ return spec.path.isBeyond(c.targetPath) }
func (spec *specification) isFirstChild() bool			{ return spec.path.lastIndex() == 0 }

func (spec *specification) execute()	{ spec.closure() }

func (spec *specification) String() string {
	return fmt.Sprintf("specification{%v @ %v}", spec.name, spec.path)
}

func asSpecArray(list *list.List) []*specification {
	arr := make([]*specification, list.Len());
	i := 0;
	for v := range list.Iter() {
		arr[i] = v.(*specification);
		i++;
	}
	return arr
}


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

