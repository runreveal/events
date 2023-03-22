//go:build !testing
// +build !testing

package text

import (
	"os"

	"github.com/runreveal/events"
)

func init() {
	handler := NewHandler("", os.Stderr)
	handler.TimeFormat = ""
	events.DefaultHandler = handler
}
