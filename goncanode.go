package goncanode

import (
	"context"
	"errors"
	"github.com/nbah1990/goncanode/api"
	"github.com/nbah1990/goncanode/entities"
	"github.com/nbah1990/goncanode/types"
)

type Handler interface {
	SignWithSecurityHeader(ctx context.Context, xml string, hashAlgorithm types.HashAlgorithm) (result entities.Response, err error)
}

func Create(o entities.Options) Handler {
	if o.Version == nil {
		v := types.NCAnodeV10
		o.Version = &v
	}

	a := api.Client{
		BaseUrl: o.ServiceUrl,
	}

	if *o.Version == types.NCAnodeV10 {
		return &NCANodeV1Handler{
			P12pass:   o.P12pass,
			P12base64: o.P12base64,
			Timeout:   o.Timeout,
			Api:       &a,
		}
	} else if *o.Version == types.NCAnodeV30 {
		return &NCANodeV3Handler{
			P12pass:   o.P12pass,
			P12base64: o.P12base64,
			Timeout:   o.Timeout,
			Api:       &a,
		}
	}

	panic(errors.New("unknown version"))
}
