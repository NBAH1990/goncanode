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

type NCANodeV1Handler struct {
	P12base64 string
	P12pass   string
	Timeout   time.Duration

	Api api.IClient
}

func (h *NCANodeV1Handler) SignWithSecurityHeader(ctx context.Context, xml string, hashAlgorithm types.HashAlgorithm) (result entities.Response, err error) {
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

func (h *NCANodeV1Handler) ExecuteRequest(ctx context.Context, r *entities.SignRequest) (result entities.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	rs, err := json.Marshal(r)
	if err != nil {
		return
	}

	rb := bytes.NewBuffer(rs)

	resp, err := h.Api.Request(ctx, http.MethodPost, ``, rb)
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
