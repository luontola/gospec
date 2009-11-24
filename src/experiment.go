// Copyright (c) 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt";
)


type Context struct{}

func (c *Context) Specify(description string, closure func()) {
	fmt.Println(description);
	closure();
}


func main() {
	count := 0;
	c := new(Context);

	c.Specify("Given the moon is full", func() {
		count++;
		c.Specify("When you walk in the woods", func() {
			count++;

			c.Specify("Then you can hear werevolves howling", func() { count++ });

			c.Specify("Then you wish you had a silver bullet", func() { count++ });
		});
	});

	c.Specify("Given the moon is not full", func() {
		count++;
		c.Specify("When you walk in the woods", func() {
			count++;

			c.Specify("Then you do not hear any werevolves", func() { count++ });

			c.Specify("Then you are not afraid", func() { count++ });
		});
	});

	fmt.Printf("number of specs: %v \n", count);
}
