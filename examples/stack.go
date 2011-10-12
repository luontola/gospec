// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

type Stack struct {
	stack []interface{}
}

func NewStack() *Stack {
	return &Stack{make([]interface{}, 0)}
}

func (this *Stack) Push(obj interface{}) {
	this.stack = append(this.stack, obj)
}

func (this *Stack) Pop() interface{} {
	last := len(this.stack) - 1
	popped := this.stack[last]
	this.stack = this.stack[:last]
	return popped
}

func (this *Stack) Empty() bool {
	return len(this.stack) == 0
}
