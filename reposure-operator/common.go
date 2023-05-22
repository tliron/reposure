package main

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
)

const toolName = "reposure-operator"

var context = contextpkg.TODO()

var log = commonlog.GetLogger(toolName)
