package main

import (
	contextpkg "context"

	"github.com/tliron/kutil/logging"
)

const toolName = "reposure-operator"

var context = contextpkg.TODO()

var log = logging.GetLogger(toolName)
