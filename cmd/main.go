package main

import (
	"fmt"
	"os"

	gitinit "github.com/cocacola-lover/kitgit/pkg/init"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Programm is run with subcommands")
		os.Exit(1)
	}

	err := runSwitch()
	if err != nil {
		fmt.Printf("Error occured : %v\n", err)
	}

}

func runSwitch() error {
	switch os.Args[1] {
	case "init":
		return gitinit.InitCmd(os.Args[2:]...)
	default:
		fmt.Println("Unknown command")
		return nil
	}
}
