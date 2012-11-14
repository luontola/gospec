// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

const (
	channelBufferSize = 10
)

// Runner executes the specs and collects their results.
type Runner struct {
	runningTasks int
	results      chan *taskResult
	executed     []*specRun
	scheduled    []*scheduledTask
	Parallel     bool
}

func NewRunner() *Runner {
	r := new(Runner)
	r.runningTasks = 0
	r.results = make(chan *taskResult, channelBufferSize)
	r.executed = make([]*specRun, 0)
	r.scheduled = make([]*scheduledTask, 0)
	r.Parallel = true
	return r
}

// Adds a spec for later execution. Example:
//     r.AddSpec(SomeSpec);
func (r *Runner) AddSpec(closure func(Context)) {
	r.AddNamedSpec(functionName(closure), closure)
}

// Adds a spec for later execution. Uses the provided name instead of
// retrieving the name of the spec function with reflection.
func (r *Runner) AddNamedSpec(name string, closure func(Context)) {
	task := newScheduledTask(name, closure, newInitialContext())
	r.scheduled = append(r.scheduled, task)
}

// Executes all the specs which have been added with AddSpec. The specs
// are executed using as many goroutines as possible, so that even individual
// spec methods are executed in multiple goroutines.
func (r *Runner) Run() {
	r.startAllScheduledTasks()
	r.startNewTasksAndWaitUntilFinished()
}

func (r *Runner) startAllScheduledTasks() {
	for r.hasScheduledTasks() {
		r.startNextScheduledTask()
	}
}

func (r *Runner) startNewTasksAndWaitUntilFinished() {
	for r.hasRunningTasks() {
		r.processNextFinishedTask()
		r.startAllScheduledTasks()
	}
}

// For testing purposes, so that the specs can be executed deterministically.
func (r *Runner) executeNextScheduledTask() {
	r.startNextScheduledTask()
	r.processNextFinishedTask()
}

func (r *Runner) startNextScheduledTask() {
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

func (r *Runner) processNextFinishedTask() {
	result := <-r.results
	r.runningTasks--
	r.saveResult(result)
}

func (r *Runner) hasRunningTasks() bool   { return r.runningTasks > 0 }
func (r *Runner) hasScheduledTasks() bool { return len(r.scheduled) > 0 }
func (r *Runner) nextScheduledTask() *scheduledTask {
	last := len(r.scheduled) - 1
	popped := r.scheduled[last]
	r.scheduled = r.scheduled[:last]
	return popped
}

func (r *Runner) execute(name string, closure specRoot, c *taskContext) *taskResult {
	c.Specify(name, func() { closure(c) })
	return &taskResult{
		name,
		closure,
		asSpecArray(c.executedSpecs),
		asSpecArray(c.postponedSpecs),
	}
}

func (r *Runner) saveResult(result *taskResult) {
	for _, spec := range result.executedSpecs {
		r.executed = append(r.executed, spec)
	}
	for _, spec := range result.postponedSpecs {
		task := newScheduledTask(result.name, result.closure, newExplicitContext(spec.path))
		r.scheduled = append(r.scheduled, task)
	}
}

func (r *Runner) Results() *ResultCollector {
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
