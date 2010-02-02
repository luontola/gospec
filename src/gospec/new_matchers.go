// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"container/list"
	"math"
	"os"
	"reflect"
)


type Matcher func(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error)

func Not(matcher Matcher) Matcher {
	return func(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
		ok, pos, neg = matcher(actual, expected)
		ok = !ok
		pos, neg = neg, pos
		return
	}
}


// The actual value must equal the expected value. For primitives the equality
// operator is used. All other objects must implement the Equality interface.
func Equals(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
	ok = areEqual(actual, expected)
	// TODO: change the messages to following?
	// '%v' should equal '%v', but it did not
	// '%v' should NOT equal '%v', but it did
	pos = Errorf("Expected '%v' but was '%v'", expected, actual)
	neg = Errorf("Did not expect '%v' but was '%v'", expected, actual)
	return
}

func areEqual(a interface{}, b interface{}) bool {
	if a2, ok := a.(Equality); ok {
		return a2.Equals(b)
	}
	return a == b
}

type Equality interface {
	Equals(other interface{}) bool
}


// TODO: IsSame - pointer equality


// The actual value must satisfy the given criteria.
func Satisfies(actual interface{}, criteria interface{}) (ok bool, pos os.Error, neg os.Error) {
	ok = criteria.(bool) == true
	pos = Errorf("Criteria not satisfied by '%v'", actual)
	neg = pos
	return
}


// The actual value must be within delta from the expected value.
func IsWithin(delta float64) Matcher {
	return func(actual_ interface{}, expected_ interface{}) (ok bool, pos os.Error, neg os.Error) {
		// TODO: need to use three-value booleans/enums for "ok" (true, false, error) or then just panic
		actual, err := toFloat64(actual_);
		if err != nil {
			return false, err, err
		}
		expected, err := toFloat64(expected_);
		if err != nil {
			return false, err, err
		}
		
		ok = math.Fabs(expected - actual) < delta
		pos = Errorf("Expected '%v' ± %v but was '%v'", expected, delta, actual)
		neg = Errorf("Did not expect '%v' ± %v but was '%v'", expected, delta, actual)
		return
	}
}

func toFloat64(actual interface{}) (result float64, err os.Error) {
	switch v := actual.(type) {
	case float:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	default:
		return math.NaN(), Errorf("Expected a float, but was '%v' of type '%T'", actual, actual)
	}
	return
}


// The actual collection (array, slice, iterator/channel) must contain the expected value.
func Contains(actual_ interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
	switch v := reflect.NewValue(actual_).(type) {
	
	case reflect.ArrayOrSliceValue:
		arr := v
		contains := false
		for i := 0; i < arr.Len(); i++ {
			other :=  arr.Elem(i).Interface()
			if areEqual(expected, other) {
				contains = true
				break
			}
		}
		ok = contains
		pos = Errorf("Expected '%v' to be in '%v' but it was not", expected, actual_)
		neg = Errorf("Did not expect '%v' to be in '%v' but it was", expected, actual_)
	
	case *reflect.ChanValue:
		ch := v
		contains := false
		list := list.New()
		for {
			other := ch.Recv().Interface()
			if ch.Closed() {
				break
			}
			if areEqual(expected, other) {
				contains = true
			}
			list.PushBack(other)
		}
		actualAsList := listToArray(list) // TODO: do lazily?
		
		ok = contains
		pos = Errorf("Expected '%v' to be in '%v' but it was not", expected, actualAsList)
		neg = Errorf("Did not expect '%v' to be in '%v' but it was", expected, actualAsList)
	
	default:
		ok = false
		pos = Errorf("Unknown type '%T', not iterable: %v", actual_, actual_)
		neg = pos
	}
	return
}

func listToArray(list *list.List) []interface{} {
	arr := make([]interface{}, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		arr[i] = e.Value
		i++
	}
	return arr
}


// TODO: ContainAll - The actual collection must contain all given elements. The order of elements is not significant.
// TODO: ContainAny - The actual collection must contain at least one element from the given collection.
// TODO: ContainExactly - The actual collection must contain exactly the same elements as in the given collection. The order of elements is not significant.
// TODO: ContainInOrder - The actual collection must contain exactly the same elements as in the given collection, and they must be in the same order.
// TODO: ContainInPartialOrder - The actual collection can hold other objects, but the objects which are common in both collections must be in the same order. The actual collection can also repeat some elements. For example [1, 2, 2, 3, 4] contains in partial order [1, 2, 3]. See Wikipedia <http://en.wikipedia.org/wiki/Partial_order> for further information.


// Helpers

func Errorf(format string, args ...) os.Error {
	return lazyString(func() string {
		return fmt.Sprintf(format, args)
	})
}

type lazyString func() string

func (this lazyString) String() string {
	return this()
}

