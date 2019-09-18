package text

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/segmentio/events/v2"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name   string
		args   bool
		output string
	}{
		{
			name: "EnableArgs:true",
			args: true,
			output: `==> 2017-01-01 23:42:00.123 - github.com/segmentio/events/text/handler_test.go:18 - Hello Luke!
	name: Luke
	from: Han
	errors:
		- EOF
`,
		},
		{
			name:   "EnableArgs:false",
			args:   false,
			output: "==> 2017-01-01 23:42:00.123 - github.com/segmentio/events/text/handler_test.go:18 - Hello Luke!\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			h := NewHandler("==> ", b)
			h.EnableArgs = test.args

			h.HandleEvent(&events.Event{
				Message: "Hello Luke!",
				Source:  "github.com/segmentio/events/text/handler_test.go:18",
				Args:    events.Args{{Name: "name", Value: "Luke"}, {Name: "from", Value: "Han"}, {Name: "error", Value: io.EOF}},
				Time:    time.Date(2017, 1, 1, 23, 42, 0, 123000000, time.Local),
				Debug:   true,
			})

			if s := b.String(); s != test.output {
				t.Error(s)
			}
		})
	}
}

func BenchmarkHandler(b *testing.B) {
	h := NewHandler("", ioutil.Discard)
	e := &events.Event{
		Message: "Hello Luke!",
		Source:  "github.com/segmentio/events/text/handler_test.go:18",
		Args:    events.Args{{Name: "name", Value: "Luke"}, {Name: "from", Value: "Han"}},
		Time:    time.Date(2017, 1, 1, 23, 42, 0, 123000000, time.UTC),
		Debug:   true,
	}

	for i := 0; i != b.N; i++ {
		h.HandleEvent(e)
	}
}
