package command

import (
	"bytes"
	"strings"
)

func (self *Client) ListImages() ([]string, error) {
	var buffer bytes.Buffer
	if err := self.Command(&buffer, "list"); err == nil {
		buffer_ := strings.TrimRight(buffer.String(), "\n")
		if buffer_ != "" {
			return strings.Split(buffer_, "\n"), nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}
