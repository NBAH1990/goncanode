package goncanode

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nbah1990/goncanode/api"
	"github.com/nbah1990/goncanode/entities"
	"github.com/nbah1990/goncanode/types"
	"net/http"
	"time"
)

type Handler interface {
	SignWithSecurityHeader(ctx context.Context, xml string, hashAlgorithm types.HashAlgorithm) (result entities.Response, err error)
	ExecuteRequest(ctx context.Context, r *entities.SignRequest) (result entities.Response, err error)
}

type NCANode struct {
	P12base64 string
	P12pass   string
	Timeout   time.Duration

	Api *api.Client
}

func Create(o entities.Options) Handler {
	a := api.Client{
		Url: o.ServiceUrl,
	}

	h := NCANode{
		P12pass:   o.P12pass,
		P12base64: o.P12base64,
		Timeout:   o.Timeout,
		Api:       &a,
	}

	return &h
}

func (h *NCANode) SignWithSecurityHeader(ctx context.Context, xml string, hashAlgorithm types.HashAlgorithm) (result entities.Response, err error) {
	r := &entities.SignRequest{
		Version:          "1.0",
		Method:           "XML.signWithSecurityHeader",
		TspHashAlgorithm: hashAlgorithm,
		Params: entities.SignParams{
			P12:      h.P12base64,
			Password: h.P12pass,
			Xml:      xml,
		},
	}

	return h.ExecuteRequest(ctx, r)
}

func (h *NCANode) ExecuteRequest(ctx context.Context, r *entities.SignRequest) (result entities.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	rs, err := json.Marshal(r)
	if err != nil {
		return
	}

	rb := bytes.NewBuffer(rs)

	resp, err := h.Api.Request(ctx, http.MethodPost, rb)
	if err != nil {
		return
	}

	var respStruct entities.Response
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return
	}

	return respStruct, nil
}
