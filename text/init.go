//go:build !testing
// +build !testing

package text

import (
	"os"

	"github.com/segmentio/events/v2"
)

func init() {
	handler := NewHandler("", os.Stderr)
	handler.TimeFormat = ""
	events.DefaultHandler = handler
}
