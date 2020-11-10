package main

import (
	"github.com/tliron/kutil/util"
)

func main() {
	err := rootCommand.Execute()
	util.FailOnError(err)
}
