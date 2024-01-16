package commands

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/go-transcribe"
)

const toolName = "reposure"

var log = commonlog.GetLogger(toolName)

var filePath string
var directoryPath string
var url string
var registry string
var tail int
var follow bool
var all bool
var sourceRegistry string
var wait bool

func Transcriber() *transcribe.Transcriber {
	return &transcribe.Transcriber{
		Strict:      strict,
		Format:      format,
		ForTerminal: pretty,
		Base64:      base64,
	}
}
