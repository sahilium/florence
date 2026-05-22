package event

import "time"

type Event struct {
	ID        string
	Source    string
	Title     string
	Body      string
	URL       string
	Severity  string
	Timestamp time.Time
	Tags      []string
}
