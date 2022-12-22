package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the little programming language!\n", user.Username)
	fmt.Println("You can type in command lines")
	fmt.Println("Happy to enjoy it")
	repl.Start(os.Stdin, os.Stdout)
}
