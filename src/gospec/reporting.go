// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt";
	"container/list";
	"sort";
)


// Collects test results for all specs in a reporting friendly format.
type ResultCollector struct {
	rootsByName map[string]*specResult;
}

func newResultCollector() *ResultCollector {
	return &ResultCollector{
		make(map[string]*specResult)
	}
}

func (r *ResultCollector) Update(spec *specRun) {
	root := r.getOrCreateRoot(spec);
	root.Update(spec);
}

func (r *ResultCollector) getOrCreateRoot(spec *specRun) *specResult {
	rawRoot := spec.rootParent();
	name := rawRoot.name;
	root, contains := r.rootsByName[name];
	if !contains {
		root = newSpecResult(rawRoot);
		r.rootsByName[name] = root;
	}
	return root
}

func (r *ResultCollector) TotalCount() int {
	totalCount := 0;
	for _, root := range r.rootsByName {
		totalCount += root.TotalCount();
	}
	return totalCount
}

func (r *ResultCollector) Roots() <-chan *specResult {
	iter := make(chan *specResult);
	go func() {
		for _, name := range r.sortedRootNames() {
			iter <- r.rootsByName[name];
		}
		close(iter);
	}();
	return iter
}

func (r *ResultCollector) sortedRootNames() []string {
	names := make([]string, len(r.rootsByName));
	i := 0;
	for name, _ := range r.rootsByName {
		names[i] = name;
		i++;
	}
	sort.SortStrings(names);
	return names
}


// Collects test results for one spec and its children in a reporting friendly format.
type specResult struct {
	name string;
	path path;
	children *list.List;
}

func newSpecResult(spec *specRun) *specResult {
	return &specResult{spec.name, spec.path, list.New()}
}

func (this *specResult) Update(spec *specRun) {
	// TODO: build a correct tree structure - create unseen, merge dublicates
	// TODO: update 'this' if the assert data differs
	if spec.path.isBeyond(this.path) {
		this.children.PushBack(newSpecResult(spec));
	}
}

func (this *specResult) TotalCount() int {
	return this.children.Len() + 1
}

func (this *specResult) Children() <-chan *specResult {
	iter := make(chan *specResult);
	go func() {
		for child := range this.children.Iter() {
			iter <- child.(*specResult);
		}
		close(iter);
	}();
	return iter
}

func (this *specResult) String() string {
	return fmt.Sprintf("%T{%v, %v, %d children}", this, this.name, this.path, this.children.Len());
}

