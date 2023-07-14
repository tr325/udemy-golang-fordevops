package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("This function requires command-line-args")
		os.Exit(1)
	}
	fmt.Printf("Hello world\nas.Args: %v\nArguments: %v\n", args, args[1:])
}
