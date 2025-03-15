package main

import (
	"flag"
	"fmt"
)

func main() {
	var enable bool
	name := flag.String("name", "World", "a name to say hello")
	age := flag.Int("age", 35, "age in years")
	flag.BoolVar(&enable, "e", false, "shorthand for enable")
	du := flag.Duration("d", 0, "shorthand for duration")
	flag.Parse()
	fmt.Println("Hello", *name)
	fmt.Println("Age", *age)
	fmt.Println("Enable", enable)
	fmt.Println("Duration", *du)
}
