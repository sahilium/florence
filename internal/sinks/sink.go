package sinks

import "github.com/sahilium/florence/internal/event"

type Sink interface {
	Send(event.Event) error
}
