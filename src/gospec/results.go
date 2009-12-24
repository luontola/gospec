// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"container/list"
	"sort"
)


// Collects test results for all specs in a reporting friendly format.
type ResultCollector struct {
	rootsByName map[string]*specResult
}

func newResultCollector() *ResultCollector {
	return &ResultCollector{
		make(map[string]*specResult),
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


type ResultVisitor interface {
	VisitSpec(name string, nestingLevel int, errors []string)
}

func (r *ResultCollector) Visit(visitor ResultVisitor) {
	r.visitAll(func(spec *specResult) {
		visitor.VisitSpec(spec.name, len(spec.path), listToStringArray(spec.errors))
	})
}

func listToStringArray(list *list.List) []string {
	arr := make([]string, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		arr[i] = e.Value.(string)
		i++
	}
	return arr
}


// TODO: Visit the specs only once and cache the counts? Even count them as they are added?

func (r *ResultCollector) TotalCount() int {
	count := 0
	r.visitAll(func(spec *specResult) {
		count++
	})
	return count
}

func (r *ResultCollector) PassCount() int {
	count := 0
	r.visitAll(func(spec *specResult) {
		if !spec.isFailed() {
			count++
		}
	})
	return count
}

func (r *ResultCollector) FailCount() int {
	count := 0
	r.visitAll(func(spec *specResult) {
		if spec.isFailed() {
			count++
		}
	})
	return count
}

func (r *ResultCollector) visitAll(visitor func(*specResult)) {
	for root := range r.roots() {
		root.visitAll(visitor)
	}
}

func (r *ResultCollector) roots() <-chan *specResult {
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
	for name, _ := range r.rootsByName {
		names[i] = name
		i++
	}
	sort.SortStrings(names)
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
	return &specResult{
		spec.name,
		spec.path,
		list.New(),
		spec.errors, // TODO: do a safe copy?
	}
}

func (this *specResult) isFailed() bool {
	return this.errors.Len() > 0
}

func (this *specResult) visitAll(visitor func(*specResult)) {
	visitor(this)
	for child := range this.children.Iter() {
		child.(*specResult).visitAll(visitor)
	}
}

func (this *specResult) update(spec *specRun) {
	isMe := this.path.isEqual(spec.path)
	isMyChild := this.path.isOn(spec.path) && !isMe
	isMyDirectChild := isMyChild && len(this.path)+1 == len(spec.path)

	if isMe {
		// TODO: check error messages and merge if different from previously registered
		return
	}

	if isMyDirectChild {
		if !this.isRegisteredChild(spec) {
			this.registerChild(spec)
		}
		return
	}

	if isMyChild {
		for child := range this.children.Iter() {
			child.(*specResult).update(spec)
		}
		return
	}
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

func (this *specResult) String() string {
	return fmt.Sprintf("%T{%v, %v, %d children}", this, this.name, this.path, this.children.Len())
}

