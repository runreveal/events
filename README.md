# events [![CircleCI](https://circleci.com/gh/segmentio/events.svg?style=shield)](https://circleci.com/gh/segmentio/events) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/events)](https://goreportcard.com/report/github.com/segmentio/events) [![GoDoc](https://godoc.org/github.com/segmentio/events?status.svg)](https://godoc.org/github.com/segmentio/events)
Go package for routing, formatting and publishing events produced by a program.

## Motivations

While Go's standard log package is handy it lacks crucial features, like the
ability to control the output format of the events for example. There are many
packages that provides logger implementations for Go but they often expose
complex APIs and were not designed to be efficient in terms of CPU and memory
usage.

The events package attempts to address these problems by providing high level
abstractions with highly efficient implementations. But it also goes further,
offering a new way to think about what logging is in a program, starting with
the package name, `events`, which expresses what this problem is about.
During its execution, a program produces events, and these events need to be
captured, routed, formatted and published to a persitence system in order to
be later analyzed.

The package was inspired by [this post](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)
from Dave Cheney. It borrowed a lot of the ideas but tried to find the sweet
spot between Dave's idealistic view of what logging is supposed to be, and
production constraints that we have here at Segment.

## Events

At the core of the package is the `Event` type. Instances of this type carry
the context in which the event was generated and the information related to
the event.

Events are passed from the sources that trigger them to handlers, which are
types implementing the `Handler` interface:
```go
type Handler interface {
    HandleEvent(*Event)
}
```
The sub-packages provide implementations of handlers that publish events to
various formats and locations.

## Logging

The `Logger` type is a source of events, the program uses loggers to generate
events with an API that helps the developer express its intent. Unlike a lot of
logging libraries, the logger doesn't support levels of messages, instead it
exposes a `Log` and `Debug` methods. Events generated by the `Log` method are
always produced by the logger, while those generated by `Debug` may be turned
on or off if necessary.

The package also exposes a default logger via top-level functions which cover
the needs of most programs. The `Log` and `Debug` functions support fmt-style
formatting but augment the syntax with features that make it simpler to generate
meaningful events. Refer to the package's documentation to learn more about it.

### Log message formatting

The `events` package supports a superset of the `fmt` formatting language. The
percent-base notation for placeholders is enhanced to automatically generated
event arguments from values passed to the call to `Log` or `Debug` functions.
This works by inserting an argument name wrapped in braces (`{}`) between the
`%` sign and the verb of the format.

For example, this piece of code generates an event that has an argument named
"name" and a value named "Luke":
```go
package main

import (
    "github.com/segmentio/events/v2"
)

func main() {
    events.Log("Hello %{name}s!", "Luke")
}
```

Note that using the extended syntax is optional and the regular `fmt` format
is supported as well.

### Compatibility with the standard library

The standard `log` package doesn't give much flexibility when it comes to its
logger type. It is a concrete type and there is no `Logger` interface which
would make it easy to plugin different implementations in packages that need to
log events. Unfortunately many of these packages have hard dependencies on the
standard logger, making it hard to capture their events and produce them in
different formats.
However, the `events/log` package is a shim between the standard `log` package,
and a stream of events. It exposes an API compatible with the standard library,
and automatically configures the `log` package to reroute the messages it emits
as events to the default logger.

## Handlers

Event handlers are the abstraction layer that allows to connect event sources to
arbitrary processing pipelines.
The sub-packages provides pre-defiend implementations of handlers.

### text

The `events/text` package provides the implementation of an event handler which
formats the event it receives in a human-readable format.

### ecs-logs

The `events/ecslogs` package provides the implementation of an event handler
which formats the events it receives in a format that is understood by ecs-logs.

We said the logger doesn't support log levels, however these levels have proven
useful to get a signal on a program misbehaving when it starts emitting tons of
*ERROR* level messages.
However, the program doesn't have to express what the severity level is in order
to get the right behavior. The `events/ecslogs` package analyzes the events it
receives and guess what the level should be, here are the rules:
- By default events are set to the *INFO* level.
- If an event was generated from a `Debug` call then handler sets the event
level to *DEBUG*.
- If the event's arguments contains at least one value that satisfies the
`error` interface then the level is set to *ERROR*.
These rules allow for the best of both worlds, giving the program a small and
expressive API to produce events while maintaining compatibility with our
existing tools.

#### DEBUG/INFO/ERROR

The events package has two main log levels (`events.Log` and `events.Debug`),
but the `ecslogs` subpackage will automatically extract error values in the
event arguments, generate _ERROR_ level messages, and put the error and stack
trace (if any is available) into the event data.

For example, this code will output a structured log message with _ERROR_ level.
```go
package main

import (
    "errors"
    "os"

    "github.com/segmentio/events/v2"
    "github.com/segmentio/events/v2/ecslogs"
)

func main() {
    events.DefaultHandler = ecslogs.NewHandler(os.Stdout)
    events.Log("something went wrong: %{error}v", errors.New("oops!"))
}
```

Otherwise, events generated by a call to `Log` will be shown as _INFO_ messages
and events generated by a call to `Debug` will be shown as _DEBUG_ messages.

### Automatic Configuration

The sub-packages have side-effects when they are imported:

Both `events/text` and `events/ecslogs` override the default logger's handler
when you import them. If the program's output is a terminal, `events/text` will
set a handler, while `events/ecslogs` will set a handler if the output is _not_
a terminal.

This approach mimics what we've achieved in many parts of our software stack and
has proven to be a good default, doing the right thing whether the program is
dealing with a production or development environment.

Here's a code example that is commonly used to configure the events package:

```go
package main

import (
    "github.com/segmentio/events/v2"
    _ "github.com/segmentio/events/v2/ecslogs"
    _ "github.com/segmentio/events/v2/text"
)

func main() {
    events.Log("enjoy!")
}
```

### Errata

This package only officially supports the latest two Go versions.
