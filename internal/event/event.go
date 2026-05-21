package event

import "time"

type Event struct {
	Source    string
	Title     string
	Body      string
	URL       string
	Severity  string
	Timestamp time.Time
	Tags      []string
}
