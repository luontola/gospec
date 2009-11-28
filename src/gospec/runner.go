// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector";
)


// SpecRunner executes the specs and collects their results.
type SpecRunner struct {
	runningTasks int;
	results chan *taskResult;
	executed *vector.Vector;
	scheduled *vector.Vector;
}

func NewSpecRunner() *SpecRunner {
	r := new(SpecRunner);
	r.runningTasks = 0;
	r.results = make(chan *taskResult);
	r.executed = vector.New(0);
	r.scheduled = vector.New(0);
	return r
}

// Adds a spec for later execution. The name of the spec method must be provided,
// because the program does not know how to find it out at runtime. Example:
//     r.AddSpec("SomeSpec", SomeSpec);
func (r *SpecRunner) AddSpec(name string, closure func(*Context)) {
	r.scheduled.Push(newScheduledTask(name, closure, newInitialContext()));
}

// Executes all the specs which have been added with AddSpec. The specs
// are executed using as many goroutines as possible, so that even individual
// spec methods are executed in multiple goroutines.
func (r *SpecRunner) Run() {
	r.startAllScheduledTasks();
	r.startNewTasksAndWaitUntilFinished();
}

func (r *SpecRunner) startAllScheduledTasks() {
	for r.hasScheduledTasks() {
		r.startNextScheduledTask();
	}
}

func (r *SpecRunner) startNewTasksAndWaitUntilFinished() {
	for r.hasRunningTasks() {
		r.processNextFinishedTask();
		r.startAllScheduledTasks();
	}
}

// For testing purposes, so that the specs can be executed deterministically.
func (r *SpecRunner) executeNextScheduledTaskSingleThreadedly() {
	r.startNextScheduledTask();
	r.processNextFinishedTask();
}

func (r *SpecRunner) startNextScheduledTask() {
	task := r.nextScheduledTask();
	go func() {
		r.results <- r.execute(task.name, task.closure, task.context);
	}();
	r.runningTasks++;
}

func (r *SpecRunner) processNextFinishedTask() {
	result := <-r.results;
	r.runningTasks--;
	r.saveResult(result);
}

func (r *SpecRunner) hasRunningTasks() bool		{ return r.runningTasks > 0 }
func (r *SpecRunner) hasScheduledTasks() bool		{ return r.scheduled.Len() > 0 }
func (r *SpecRunner) nextScheduledTask() *scheduledTask	{ return r.scheduled.Pop().(*scheduledTask) }

func (r *SpecRunner) execute(name string, closure func(*Context), c *Context) *taskResult {
	c.Specify(name, func() { closure(c) });
	return &taskResult{
		name, closure,
		asSpecArray(c.executedSpecs),
		asSpecArray(c.postponedSpecs)
	}
}

func (r *SpecRunner) saveResult(result *taskResult) {
	for _, spec := range result.executedSpecs {
		r.executed.Push(spec);
	}
	for _, spec := range result.postponedSpecs {
		task := newScheduledTask(result.name, result.closure, newExplicitContext(spec.path));
		r.scheduled.Push(task);
	}
}


// Scheduled spec execution.
type scheduledTask struct {
	name string;
	closure func(*Context);
	context *Context;
}

func newScheduledTask(name string, closure func(*Context), context *Context) *scheduledTask {
	return &scheduledTask{name, closure, context}
}


// Results of a spec execution.
type taskResult struct {
	name string;
	closure func(*Context);
	executedSpecs []*specification;
	postponedSpecs []*specification;
}

