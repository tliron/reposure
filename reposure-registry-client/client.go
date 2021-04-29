package main

import (
	"github.com/tliron/reposure/client/direct"
)

func NewClient() *direct.Client {
	return direct.NewClient(host, roundTripper, username, password, token, context)
}
