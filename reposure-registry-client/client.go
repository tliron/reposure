package main

import (
	directclient "github.com/tliron/reposure/client/direct"
)

func NewClient() *directclient.Client {
	return directclient.NewClient(host, roundTripper, username, password, token, context)
}
