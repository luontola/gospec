// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector";
)


// Starts the spec execution and collects the results

type SpecRunner struct {
	executed *vector.Vector;
	scheduled *vector.Vector;
}

func NewSpecRunner() *SpecRunner {
	r := new(SpecRunner);
	r.executed = vector.New(0);
	r.scheduled = vector.New(0);
	return r
}

func (r *SpecRunner) AddSpec(name string, closure func(*Context)) {
	r.scheduled.Push(newScheduledTask(name, closure, newInitialContext()));
}

func (r *SpecRunner) Run() {
	for r.hasScheduledTasks() {
		r.executeNextScheduledTask();
	}
}

func (r *SpecRunner) executeNextScheduledTask() {
	task := r.nextScheduledTask();
	result := r.execute(task.name, task.closure, task.context);
	r.saveResult(task, result);
}

func (r *SpecRunner) hasScheduledTasks() bool		{ return r.scheduled.Len() > 0 }
func (r *SpecRunner) nextScheduledTask() *scheduledTask	{ return r.scheduled.Pop().(*scheduledTask) }

func (r *SpecRunner) execute(name string, closure func(*Context), c *Context) *runResult {
	c.Specify(name, func() { closure(c) });
	return &runResult{
		asSpecArray(c.executedSpecs),
		asSpecArray(c.postponedSpecs)
	}
}

func (r *SpecRunner) saveResult(task *scheduledTask, result *runResult) {
	for _, spec := range result.executedSpecs {
		r.executed.Push(spec);
	}
	for _, spec := range result.postponedSpecs {
		r.scheduled.Push(task.copy(spec.path));
	}
}


// Scheduled spec execution

type scheduledTask struct {
	name string;
	closure func(*Context);
	context *Context;
}

func newScheduledTask(name string, closure func(*Context), context *Context) *scheduledTask {
	return &scheduledTask{name, closure, context}
}

func (task *scheduledTask) copy(targetPath path) *scheduledTask {
	return newScheduledTask(task.name, task.closure, newExplicitContext(targetPath))
}


// Results of a spec execution

type runResult struct {
	executedSpecs []*specification;
	postponedSpecs []*specification;
}

