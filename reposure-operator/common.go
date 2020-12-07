package main

import (
	contextpkg "context"

	"github.com/op/go-logging"
)

const toolName = "reposure-operator"

var context = contextpkg.TODO()

var log = logging.MustGetLogger(toolName)
