package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/util"
)

func main() {
	err := rootCommand.Execute()
	util.FailOnError(err)
	atexit.Exit(0)
}
