// Copyright (c) 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt";
)


type Context struct {
	level int;
}

func (c *Context) Specify(description string, closure func()) {
	for i := 0; i < c.level; i++ {
		fmt.Print("  ")
	}
	fmt.Println(" - " + description);
	c.level++;
	closure();
	c.level--;
}


func main() {
	count := 0;
	c := new(Context);

	/*
	  GoSpec will execute the following specs in this order, so that each
	  of these rows is executed in its own isolated goroutine:
	  - 1, 1.1, 1.1.1
	  - 1, 1.1, 1.1.2
	  - 2, 2.1, 2.1.1
	  - 2, 2.1, 2.1.2
	*/

	c.Specify("Given the moon is full", func() {
		count++;	// 1
		c.Specify("When you walk in the woods", func() {
			count++;	// 1.1

			c.Specify("Then you can hear werevolves howling", func() {
				count++	// 1.1.1
			});

			c.Specify("Then you wish you had a silver bullet", func() {
				count++	// 1.1.2
			});
		});
	});

	c.Specify("Given the moon is not full", func() {
		count++;	// 2
		c.Specify("When you walk in the woods", func() {
			count++;	// 2.1

			c.Specify("Then you do not hear any werevolves", func() {
				count++	// 2.1.1
			});

			c.Specify("Then you are not afraid", func() {
				count++	// 2.1.2
			});
		});
	});

	fmt.Printf("number of specs: %v \n", count);
}
