// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

const (
	channelBufferSize = 10
)

type Runner interface {
	Run()

	AddSpec(func(Context))
	AddNamedSpec(string, func(Context))

	Results() *ResultCollector
}

// Runner executes the specs and collects their results.
type ParallelRunner struct {
	runningTasks int
	results      chan *taskResult
	executed     []*specRun
	scheduled    []*scheduledTask
	Parallel     bool
}

// TODO Implement this as a type
func NewSerialRunner() *ParallelRunner {
	r := new(ParallelRunner)
	r.runningTasks = 0
	r.results = make(chan *taskResult, channelBufferSize)
	r.executed = make([]*specRun, 0)
	r.scheduled = make([]*scheduledTask, 0)
	return r
}

func NewParallelRunner() *ParallelRunner {
	r := NewSerialRunner()
	r.Parallel = true
	return r
}

// Adds a spec for later execution. Example:
//     r.AddSpec(SomeSpec);
func (r *ParallelRunner) AddSpec(closure func(Context)) {
	r.AddNamedSpec(functionName(closure), closure)
}

// Adds a spec for later execution. Uses the provided name instead of
// retrieving the name of the spec function with reflection.
func (r *ParallelRunner) AddNamedSpec(name string, closure func(Context)) {
	task := newScheduledTask(name, closure, newInitialContext())
	r.scheduled = append(r.scheduled, task)
}

// Executes all the specs which have been added with AddSpec. The specs
// are executed using as many goroutines as possible, so that even individual
// spec methods are executed in multiple goroutines.
func (r *ParallelRunner) Run() {
	r.startAllScheduledTasks()
	r.startNewTasksAndWaitUntilFinished()
}

func (r *ParallelRunner) startAllScheduledTasks() {
	for r.hasScheduledTasks() {
		r.startNextScheduledTask()
	}
}

func (r *ParallelRunner) startNewTasksAndWaitUntilFinished() {
	for r.hasRunningTasks() {
		r.processNextFinishedTask()
		r.startAllScheduledTasks()
	}
}

// For testing purposes, so that the specs can be executed deterministically.
func (r *ParallelRunner) executeNextScheduledTask() {
	r.startNextScheduledTask()
	r.processNextFinishedTask()
}

func (r *ParallelRunner) startNextScheduledTask() {
	task := r.nextScheduledTask()
	if r.Parallel {
		go func() {
			r.results <- r.execute(task.name, task.closure, task.context)
		}()
	} else {
		r.results <- r.execute(task.name, task.closure, task.context)
	}
	r.runningTasks++
}

func (r *ParallelRunner) processNextFinishedTask() {
	result := <-r.results
	r.runningTasks--
	r.saveResult(result)
}

func (r *ParallelRunner) hasRunningTasks() bool   { return r.runningTasks > 0 }
func (r *ParallelRunner) hasScheduledTasks() bool { return len(r.scheduled) > 0 }
func (r *ParallelRunner) nextScheduledTask() *scheduledTask {
	last := len(r.scheduled) - 1
	popped := r.scheduled[last]
	r.scheduled = r.scheduled[:last]
	return popped
}

func (r *ParallelRunner) execute(name string, closure specRoot, c *taskContext) *taskResult {
	c.Specify(name, func() { closure(c) })
	return &taskResult{
		name,
		closure,
		asSpecArray(c.executedSpecs),
		asSpecArray(c.postponedSpecs),
	}
}

func (r *ParallelRunner) saveResult(result *taskResult) {
	for _, spec := range result.executedSpecs {
		r.executed = append(r.executed, spec)
	}
	for _, spec := range result.postponedSpecs {
		task := newScheduledTask(result.name, result.closure, newExplicitContext(spec.path))
		r.scheduled = append(r.scheduled, task)
	}
}

func (r *ParallelRunner) Results() *ResultCollector {
	// TODO: Should this be done concurrently with executing the specs?
	// The result collector could run in its own goroutine, and the
	// Runner.saveResult() method would send executed specs to it as they
	// get ready (the channel should be buffered). When all is done, the runner
	// will get the result collector from a result channel.

	results := newResultCollector()
	for _, spec := range r.executed {
		results.Update(spec)
	}
	return results
}

// Scheduled spec execution.
type scheduledTask struct {
	name    string
	closure specRoot
	context *taskContext
}

type specRoot func(Context)

func newScheduledTask(name string, closure specRoot, context *taskContext) *scheduledTask {
	return &scheduledTask{name, closure, context}
}

// Results of a spec execution.
type taskResult struct {
	name           string
	closure        specRoot
	executedSpecs  []*specRun
	postponedSpecs []*specRun
}
