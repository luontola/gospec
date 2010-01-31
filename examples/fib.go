// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples


type Fib struct {
	n0, n1 int
}

func NewFib() *Fib {
	return &Fib{0, 1}
}

func (this *Fib) Next() int {
	next := this.n0
	this.n0, this.n1 = this.n1, this.n0+this.n1
	return next
}

func (this *Fib) Sequence(length int) []int {
	seq := make([]int, length)
	for i := range seq {
		seq[i] = this.Next()
	}
	return seq
}

