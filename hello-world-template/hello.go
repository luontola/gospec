package hello

import (
	"fmt"
)


func SayHello(name string) string {
	return fmt.Sprintf("Hello, %v", name)
}
