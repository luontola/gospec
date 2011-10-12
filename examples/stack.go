// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"container/vector"
)

type Stack struct {
	stack *vector.Vector
}

func NewStack() *Stack {
	return &Stack{new(vector.Vector)}
}

func (this *Stack) Push(obj interface{}) {
	this.stack.Push(obj)
}

func (this *Stack) Pop() interface{} {
	return this.stack.Pop()
}

func (this *Stack) Empty() bool {
	return this.stack.Len() == 0
}
