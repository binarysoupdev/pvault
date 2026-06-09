package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/command/sample"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	runner := command.NewRunner(sample.NewHelloCommand())

	if *ls || len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}
