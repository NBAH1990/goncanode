package entities

import "github.com/nbah1990/goncanode/types"

type Response struct {
	Result  ResponseResult `json:"result"`
	Message string         `json:"message"`
	Status  int            `json:"status"`
}

type ResponseResult struct {
	Xml string `json:"xml"`
	Raw string `json:"raw"`
}

type SignRequest struct {
	Version          string              `json:"version"`
	Method           string              `json:"method"`
	TspHashAlgorithm types.HashAlgorithm `json:"tspHashAlgorithm"`
	Params           SignParams          `json:"params"`
}

type SignParams struct {
	P12      string `json:"p12"`
	Password string `json:"password"`
	Xml      string `json:"xml"`
}
