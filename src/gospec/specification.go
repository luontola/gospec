// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"fmt"
)


// Represents a spec in a tree of specs.
type specRun struct {
	name             string
	closure          func()
	parent           *specRun
	numberOfChildren int
	path             path
	errors           *list.List
	hasFatalErrors   bool
}

func newSpecRun(name string, closure func(), parent *specRun) *specRun {
	path := rootPath()
	if parent != nil {
		currentIndex := parent.numberOfChildren
		path = parent.path.append(currentIndex)
		parent.numberOfChildren++
	}
	return &specRun{name, closure, parent, 0, path, list.New(), false}
}

func (spec *specRun) isOnTargetPath(c *Context) bool { return spec.path.isOn(c.targetPath) }
func (spec *specRun) isUnseen(c *Context) bool       { return spec.path.isBeyond(c.targetPath) }
func (spec *specRun) isFirstChild() bool             { return spec.path.lastIndex() == 0 }

func (spec *specRun) execute() {
	spec.closure()
}

func (spec *specRun) AddError(message string) {
	spec.errors.PushBack(message)
}

func (spec *specRun) AddFatalError(message string) {
	spec.AddError(message)
	spec.hasFatalErrors = true
}

func (spec *specRun) rootParent() *specRun {
	root := spec
	for root.parent != nil {
		root = root.parent
	}
	return root
}

func (spec *specRun) String() string {
	return fmt.Sprintf("%T{%v @ %v}", spec, spec.name, spec.path)
}


func asSpecArray(list *list.List) []*specRun {
	arr := make([]*specRun, list.Len())
	i := 0
	for v := range list.Iter() {
		arr[i] = v.(*specRun)
		i++
	}
	return arr
}


// Path of a specification.
type path []int

func rootPath() path {
	return []int{}
}

func (parent path) append(index int) path {
	result := make([]int, len(parent)+1)
	for i, v := range parent {
		result[i] = v
	}
	result[len(parent)] = index
	return result
}

func (current path) isOn(target path) bool {
	return commonPrefixLength(current, target) == len(current)
}

func (current path) isEqual(target path) bool {
	return current.isOn(target) && len(current) == len(target)
}

func (current path) isBeyond(target path) bool {
	return target.isOn(current) && len(current) > len(target)
}

func commonPrefixLength(a path, b path) int {
	length := 0
	for i := 0; i < len(a) && i < len(b) && a[i] == b[i]; i++ {
		length++
	}
	return length
}

func (path path) lastIndex() int {
	if len(path) == 0 {
		return -1 // root path
	}
	return path[len(path)-1]
}

func (path path) isRoot() bool {
	return len(path) == 0
}

