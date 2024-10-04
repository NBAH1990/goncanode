package entities

import (
	"github.com/nbah1990/goncanode/types"
	"time"
)

type Options struct {
	ServiceUrl string
	P12base64  string
	P12pass    string
	Timeout    time.Duration

	Version *types.Version
}
