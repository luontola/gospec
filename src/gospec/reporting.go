// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt";
	"container/list";
	"sort";
)

var x = fmt.Sprintf("keep 'fmt' import during debugging"); // TODO: remove


// Collects test results for all specs in a reporting friendly format.
type specReport struct {
	rootsByName map[string]*specInfo;
}

func newSpecReport() *specReport {
	return &specReport{make(map[string]*specInfo)}
}

func (r *specReport) Update(spec *spec) {
	root := r.getOrCreateRoot(spec);
	root.Update(spec);
}

func (r *specReport) getOrCreateRoot(spec *spec) *specInfo {
	rawRoot := spec.rootParent();
	name := rawRoot.name;
	root, contains := r.rootsByName[name];
	if !contains {
		root = newSpecInfo(rawRoot);
		r.rootsByName[name] = root;
	}
	return root
}

func (r *specReport) TotalCount() int {
	totalCount := 0;
	for _, root := range r.rootsByName {
		totalCount += root.TotalCount();
	}
	return totalCount
}

func (r *specReport) Roots() <-chan *specInfo {
	iter := make(chan *specInfo);
	go func() {
		for _, name := range r.sortedRootNames() {
			iter <- r.rootsByName[name];
		}
		close(iter);
	}();
	return iter
}

func (r *specReport) sortedRootNames() []string {
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
type specInfo struct {
	name string;
	path path;
	children *list.List;
}

func newSpecInfo(spec *spec) *specInfo {
	return &specInfo{spec.name, spec.path, list.New()}
}

func (this *specInfo) Update(spec *spec) {
	// TODO: build a correct tree structure - create unseen, merge dublicates
	// TODO: update 'this' if the assert data differs
	if spec.path.isBeyond(this.path) {
		this.children.PushBack(newSpecInfo(spec));
	}
}

func (this *specInfo) TotalCount() int {
	return this.children.Len() + 1
}

func (this *specInfo) Children() <-chan *specInfo {
	iter := make(chan *specInfo);
	go func() {
		for child := range this.children.Iter() {
			iter <- child.(*specInfo);
		}
		close(iter);
	}();
	return iter
}

func (this *specInfo) String() string {
	return fmt.Sprintf("%T{%v, %v, %d children}", this, this.name, this.path, this.children.Len());
}

