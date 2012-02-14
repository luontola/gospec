// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"container/list"
	"github.com/orfjackal/gospec/src/gospec"   // the "gospec.Context" interface
	. "github.com/orfjackal/gospec/src/gospec" // the expectation matchers (Equals, IsTrue etc.), will later be renamed to "gospec/matchers"
	"os"
)

func ExpectationSyntaxSpec(c gospec.Context) {

	c.Specify("Objects can be compared for equality", func() {
		c.Expect(1, Equals, 1)
		c.Expect("string", Equals, "string")

		// There are some shorthands for commonly used comparisons:
		c.Expect(true, IsTrue)
		c.Expect(false, IsFalse)
		c.Expect(nil, IsNil)
		var typedNilPointerInsideInterfaceValue *os.File
		c.Expect(typedNilPointerInsideInterfaceValue, IsNil)

		// Comparing pointer equality is also possible:
		p1 := &Point2{1, 2}
		p2 := p1
		p3 := &Point2{1, 2}
		c.Expect(p2, IsSame, p1)
		c.Expect(p3, Not(IsSame), p1)

		// Comparing floats for equality is not recommended, because
		// floats are rarely exactly equal. So don't write like this:
		c.Expect(3.141, Equals, 3.141)
		// But instead compare using a delta and write like this:
		c.Expect(3.141, IsWithin(0.001), 3.1415926535)

		// Objects with an "Equals(interface{}) bool" method can be
		// compared for equality. See "point.go" for details of how
		// the Equals(interface{}) method should be written. Special
		// care is needed if the objects are used both as values and
		// as pointers.
		a1 := Point2{1, 2}
		a2 := Point2{1, 2}
		c.Expect(a1, Equals, a2)

		b1 := &Point3{1, 2, 3}
		b2 := &Point3{1, 2, 3}
		c.Expect(b1, Equals, b2)
	})

	c.Specify("All expectations can be negated", func() {
		c.Expect(1, Not(Equals), 2)
		c.Expect("apples", Not(Equals), "oranges")
		c.Expect(new(int), Not(IsNil))
	})

	c.Specify("Boolean expressions can be stated about an object", func() {
		s := "some string"
		c.Expect(s, Satisfies, len(s) >= 10 && len(s) <= 20)
		c.Expect(s, Not(Satisfies), len(s) == 0)
	})

	c.Specify("Custom matchers can be defined for commonly used expressions", func() {
		c.Expect("first string", HasSameLengthAs, "other string")
	})

	c.Specify("Arrays/slices, lists and channels can be tested for containment", func() {
		array := []string{"one", "two", "three"}
		list := list.New()
		list.PushBack("one")
		list.PushBack("two")
		list.PushBack("three")
		channel := make(chan string, 10)
		channel <- "one"
		channel <- "two"
		channel <- "three"
		close(channel)

		c.Expect(array, Contains, "one")
		c.Expect(list, Contains, "two")
		c.Expect(channel, Contains, "three")
		c.Expect(array, Not(Contains), "four")

		c.Expect(list, ContainsAll, Values("two", "one"))
		c.Expect(list, ContainsAny, Values("apple", "orange", "one"))
		c.Expect(list, ContainsExactly, Values("two", "one", "three"))
		c.Expect(list, ContainsInOrder, Values("one", "two", "three"))
		c.Expect(list, ContainsInPartialOrder, Values("one", "three"))
	})
}

func HasSameLengthAs(actual interface{}, expected interface{}) (match bool, pos gospec.Message, neg gospec.Message, err error) {
	lenActual := len(actual.(string))
	lenExpected := len(expected.(string))
	difference := lenActual - lenExpected

	match = difference == 0
	pos = gospec.Messagef(actual, "has same length as “%v” (difference was %+d)", expected, difference)
	neg = gospec.Messagef(actual, "does NOT have same length as “%v” (difference was %+d)", expected, difference)
	return
}
