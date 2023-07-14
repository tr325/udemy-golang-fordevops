package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Printf("Hello world\nas.Args: %v\nArguments: %v\n", args, args[1:])
}
