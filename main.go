package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/taylorza/zxenv/internal/engine"
)

var devpath = "."

func main() {
	if len(os.Args) < 2 {
		showUsage("")
		return
	}

	switch os.Args[1] {
	case "init":
		env := &engine.Environment{
			BasePath: filepath.ToSlash(devpath),
		}
		err := engine.SetupDevelopment(env)
		if err != nil {
			log.Fatal(err)
		}

	case "new":
		if len(os.Args) < 3 {
			showUsage("project name required")
		}
		err := engine.CreateProject(".", os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	default:
		showUsage("")
	}
}

func showUsage(message string) {
	fmt.Println("ZXEnv - CLI to install an emulation evironment and scafold new projects")
	if len(message) > 0 {
		fmt.Println()
		fmt.Println(message)
	}
	fmt.Println("Usage:")
	fmt.Println(" zxenv <cmd> [arguments]")
	fmt.Println()
	fmt.Println(" Commands:")
	fmt.Println("  init      \t- Initializes the current directory with a full development environment")
	fmt.Println("  new <name>\t- Creates a new template project with specified name")
	os.Exit(1)
}
