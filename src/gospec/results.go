// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"fmt"
	"sort"
)

// Collects test results for all specs in a reporting friendly format.
type ResultCollector struct {
	rootsByName map[string]*specResult
	passCount   int
	failCount   int
}

func newResultCollector() *ResultCollector {
	return &ResultCollector{
		make(map[string]*specResult),
		-1,
		-1,
	}
}

func (r *ResultCollector) Update(spec *specRun) {
	root := r.getOrCreateRoot(spec)
	root.update(spec)
}

func (r *ResultCollector) getOrCreateRoot(spec *specRun) *specResult {
	rawRoot := spec.rootParent()
	name := rawRoot.name
	root, contains := r.rootsByName[name]
	if !contains {
		root = newSpecResult(rawRoot)
		r.rootsByName[name] = root
	}
	return root
}

// Number of specs

func (r *ResultCollector) TotalCount() int {
	return r.PassCount() + r.FailCount()
}

func (r *ResultCollector) PassCount() int {
	if r.passCount < 0 {
		r.calculateSpecCount()
	}
	return r.passCount
}

func (r *ResultCollector) FailCount() int {
	if r.failCount < 0 {
		r.calculateSpecCount()
	}
	return r.failCount
}

func (r *ResultCollector) calculateSpecCount() {
	r.resetSpecCount()
	r.visitAll(func(spec *specResult) {
		r.incrementSpecCount(spec)
	})
}

func (r *ResultCollector) resetSpecCount() {
	r.failCount = 0
	r.passCount = 0
}

func (r *ResultCollector) incrementSpecCount(spec *specResult) {
	if spec.isFailed() {
		r.failCount++
	} else {
		r.passCount++
	}
}

// Visiting the results

type ResultVisitor interface {
	VisitSpec(nestingLevel int, name string, errors []*Error)
	VisitEnd(passCount int, failCount int)
}

func (r *ResultCollector) Visit(visitor ResultVisitor) {
	r.resetSpecCount()
	r.visitAll(func(spec *specResult) {
		r.incrementSpecCount(spec)
		visitor.VisitSpec(len(spec.path), spec.name, listToErrorArray(spec.errors))
	})
	visitor.VisitEnd(r.passCount, r.failCount)
}

func listToErrorArray(list *list.List) []*Error {
	arr := make([]*Error, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		arr[i] = e.Value.(*Error)
		i++
	}
	return arr
}

func (r *ResultCollector) visitAll(visitor func(*specResult)) {
	for root := range r.sortedRoots() {
		root.visitAll(visitor)
	}
}

func (r *ResultCollector) sortedRoots() <-chan *specResult {
	iter := make(chan *specResult)
	go func() {
		for _, name := range r.sortedRootNames() {
			iter <- r.rootsByName[name]
		}
		close(iter)
	}()
	return iter
}

func (r *ResultCollector) sortedRootNames() []string {
	names := make([]string, len(r.rootsByName))
	i := 0
	for name := range r.rootsByName {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

// Collects test results for one spec and its children in a reporting friendly format.
type specResult struct {
	name     string
	path     path
	children *list.List
	errors   *list.List
}

func newSpecResult(spec *specRun) *specResult {
	// 'children' and 'errors' will be populated by update()
	return &specResult{
		spec.name,
		spec.path,
		list.New(),
		list.New(),
	}
}

func (this *specResult) isFailed() bool {
	return this.errors.Len() > 0
}

func (this *specResult) visitAll(visitor func(*specResult)) {
	visitor(this)
	for e := this.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*specResult)
		child.visitAll(visitor)
	}
}

func (this *specResult) update(spec *specRun) {
	isMe := this.path.isEqual(spec.path)
	isMyChild := this.path.isOn(spec.path) && !isMe
	isMyDirectChild := isMyChild && len(this.path)+1 == len(spec.path)

	if isMe {
		this.mergeErrors(spec.errors)
	}
	if isMyDirectChild {
		if !this.isRegisteredChild(spec) {
			this.registerChild(spec)
		}
	}
	if isMyChild {
		child := this.findChildOnPath(spec.path)
		child.update(spec)
	}
}

func (this *specResult) mergeErrors(newErrors *list.List) {
	for e := newErrors.Front(); e != nil; e = e.Next() {
		error := e.Value.(*Error)
		if !this.hasError(error) {
			this.addError(error)
		}
	}
}

func (this *specResult) hasError(error *Error) bool {
	for e := this.errors.Front(); e != nil; e = e.Next() {
		if error.equals(e.Value.(*Error)) {
			return true
		}
	}
	return false
}

func (this *specResult) addError(error *Error) {
	this.errors.PushBack(error)
}

func (this *specResult) isRegisteredChild(spec *specRun) bool {
	for e := this.children.Front(); e != nil; e = e.Next() {
		other := e.Value.(*specResult)
		if other.path.isEqual(spec.path) {
			return true
		}
	}
	return false
}

func (this *specResult) registerChild(spec *specRun) {
	newChild := newSpecResult(spec)
	pos := this.findFirstChildWithGreaterIndex(newChild.path.lastIndex())
	if pos != nil {
		this.children.InsertBefore(newChild, pos)
	} else {
		this.children.PushBack(newChild)
	}
}

func (this *specResult) findFirstChildWithGreaterIndex(targetIndex int) *list.Element {
	for e := this.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*specResult)
		if child.path.lastIndex() > targetIndex {
			return e
		}
	}
	return nil
}

func (this *specResult) findChildOnPath(targetPath path) *specResult {
	for e := this.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*specResult)
		if child.path.isOn(targetPath) {
			return child
		}
	}
	return nil
}

func (this *specResult) String() string {
	return fmt.Sprintf("%T{%v, %v, %d children, %d errors}",
		this, this.name, this.path, this.children.Len(), this.errors.Len())
}
