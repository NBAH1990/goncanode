package entities

import "time"

type Options struct {
	ServiceUrl string
	P12base64  string
	P12pass    string
	Timeout    time.Duration
}
