package main

import (
	"github.com/tliron/kutil/util"
	"github.com/tliron/reposure/reposure/commands"

	_ "github.com/tliron/commonlog/simple"
)

func main() {
	util.ExitOnSignals()
	commands.Execute()
	util.Exit(0)
}
