package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/reposure/reposure/commands"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}
