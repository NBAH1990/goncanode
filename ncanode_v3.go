package goncanode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbah1990/goncanode/api"
	"github.com/nbah1990/goncanode/entities"
	"github.com/nbah1990/goncanode/types"
	"net/http"
	"time"
)

type NCANodeV3Handler struct {
	P12base64 string
	P12pass   string
	Timeout   time.Duration

	Api api.IClient
}

type wsseSignRequest struct {
	Xml      string  `json:"xml"`
	Key      string  `json:"key"`
	Password string  `json:"password"`
	KeyAlias *string `json:"keyAlias"`
	TrimXml  bool    `json:"trimXml"`
}

type wsseSignResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Xml     string `json:"xml"`
}

func (h *NCANodeV3Handler) SignWithSecurityHeader(ctx context.Context, xmlS string, _ types.HashAlgorithm) (result entities.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	r := wsseSignRequest{
		Xml:      xmlS,
		Key:      h.P12base64,
		Password: h.P12pass,
		TrimXml:  false,
		KeyAlias: nil,
	}

	rs, err := json.Marshal(r)
	if err != nil {
		return result, errors.New(fmt.Sprintf(`SignXml: can't encode request json: %s`, err))
	}

	rb := bytes.NewBuffer(rs)

	resp, err := h.Api.Request(ctx, http.MethodPost, `/wsse/sign`, rb)
	if err != nil {
		return result, errors.New(fmt.Sprintf(`SignXml: http request error: %s`, err))
	}

	var respStruct wsseSignResponse
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return result, errors.New(fmt.Sprintf(`SignXml: can't decode http response json: %s`, err))
	}

	if respStruct.Status != http.StatusOK {
		return result, errors.New(fmt.Sprintf(`SignXml: http error: %s, status: %d`, respStruct.Message, respStruct.Status))
	}

	result.Result.Xml = respStruct.Xml
	result.Result.Raw = respStruct.Xml
	result.Status = respStruct.Status
	result.Message = respStruct.Message

	return result, nil
}
