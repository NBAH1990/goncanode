package goncanode

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/nbah1990/goncanode/entities"
	"github.com/nbah1990/goncanode/types"
)

type mockApiClientV1 struct {
	response []byte
	err      error
}

func (m *mockApiClientV1) Request(_ context.Context, _ string, _ string, _ *bytes.Buffer) ([]byte, error) {
	return m.response, m.err
}

func TestNCANodeV1Handler_SignWithSecurityHeader(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler := &NCANodeV1Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClientV1{
				response: []byte(`{"status":200,"message":"Success","result":{"xml":"<signedXml></signedXml>"}}`),
				err:      nil,
			},
		}

		ctx := context.Background()
		result, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", types.GOST34311GT)
		if err != nil {
			t.Fatalf("Expected no errors, got: %v", err)
		}
		if result.Result.Xml != "<signedXml></signedXml>" {
			t.Errorf("Incorrect Xml result, got: %s", result.Result.Xml)
		}
	})

	t.Run("JsonMarshalError", func(t *testing.T) {
		handler := &NCANodeV1Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api:       &mockApiClientV1{},
		}

		ctx := context.Background()
		// Force a marshalling error by setting invalid Params
		invalidParams := entities.SignRequest{
			Version: "1.0",
			Method:  "XML.signWithSecurityHeader",
			Params:  entities.SignParams{},
		}
		invalidParams.Params.Xml = string([]byte{0xff, 0xfe, 0xfd})

		_, err := handler.ExecuteRequest(ctx, &invalidParams)
		if err == nil {
			t.Errorf("Expected JSON marshalling error, got nil")
		}
	})

	t.Run("ApiRequestError", func(t *testing.T) {
		handler := &NCANodeV1Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClientV1{
				err: errors.New("request error"),
			},
		}

		ctx := context.Background()
		_, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", types.GOST34311GT)
		if err == nil || err.Error() != "request error" {
			t.Errorf("Expected API request error, got: %v", err)
		}
	})

	t.Run("JsonUnmarshalError", func(t *testing.T) {
		handler := &NCANodeV1Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClientV1{
				response: []byte(`invalid json`),
			},
		}

		ctx := context.Background()
		_, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", types.GOST34311GT)
		if err == nil {
			t.Errorf("Expected JSON unmarshalling error, got nil")
		}
	})

	t.Run("NonOKStatus", func(t *testing.T) {
		handler := &NCANodeV1Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClientV1{
				response: []byte(`{"status":400,"message":"Bad Request","result":{}}`),
			},
		}

		ctx := context.Background()
		result, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", types.GOST34311GT)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result.Status != 400 {
			t.Errorf("Expected status 400, got: %d", result.Status)
		}
		if result.Message != "Bad Request" {
			t.Errorf("Expected message 'Bad Request', got: %s", result.Message)
		}
	})
}
