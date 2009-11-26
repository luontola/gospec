
package gospec

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



