package commands

import (
	"github.com/op/go-logging"
)

const toolName = "reposure"

var log = logging.MustGetLogger(toolName)

var filePath string
var directoryPath string
var url string
var repository string
var tail int
var follow bool
var all bool
var registry string
var wait bool
