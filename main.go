package main

import (
	"flag"
	"fmt"
)

var (
	 usage = "this is the usage of the argument"
)

func main() {
	fmt.Println("test one and two ")
	var argumentName = flag.String("osintmaster", "valeu", usage)
	fmt.Println(string(*argumentName))
}
