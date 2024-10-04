package goncanode

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

type mockApiClient struct {
	response []byte
	err      error
}

func (m *mockApiClient) Request(_ context.Context, _ string, _ string, _ *bytes.Buffer) ([]byte, error) {
	return m.response, m.err
}

func TestSignXml(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler := &NCANodeV3Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClient{
				response: []byte(`{"status":200,"message":"Success","xml":"<signedXml></signedXml>"}`),
				err:      nil,
			},
		}

		ctx := context.Background()
		result, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", ``)
		if err != nil {
			t.Fatalf("Expected no errors, got: %v", err)
		}
		if result.Result.Xml != "<signedXml></signedXml>" {
			t.Errorf("Incorrect Xml result, got: %s", result.Result.Xml)
		}
	})

	t.Run("JsonMarshalError", func(t *testing.T) {
		handler := &NCANodeV3Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api:       &mockApiClient{},
		}

		ctx := context.Background()
		// Force a marshalling error by passing invalid xmlS
		_, err := handler.SignWithSecurityHeader(ctx, string([]byte{0xff, 0xfe, 0xfd}), ``)
		if err == nil || err.Error() != `SignXml: can't decode http response json: unexpected end of JSON input` {
			t.Errorf("Expected JSON marshalling error, got: %v", err)
		}
	})

	t.Run("ApiRequestError", func(t *testing.T) {
		handler := &NCANodeV3Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClient{
				err: errors.New("request error"),
			},
		}

		ctx := context.Background()
		_, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", ``)
		if err == nil || err.Error() != "SignXml: http request error: request error" {
			t.Errorf("Expected API request error, got: %v", err)
		}
	})

	t.Run("JsonUnmarshalError", func(t *testing.T) {
		handler := &NCANodeV3Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClient{
				response: []byte(`invalid json`),
			},
		}

		ctx := context.Background()
		_, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", ``)
		if err == nil || err.Error() != "SignXml: can't decode http response json: invalid character 'i' looking for beginning of value" {
			t.Errorf("Expected JSON unmarshalling error, got: %v", err)
		}
	})

	t.Run("NonOKStatus", func(t *testing.T) {
		handler := &NCANodeV3Handler{
			P12base64: "base64string",
			P12pass:   "password",
			Api: &mockApiClient{
				response: []byte(`{"status":400,"message":"Bad Request","xml":""}`),
			},
		}

		ctx := context.Background()
		_, err := handler.SignWithSecurityHeader(ctx, "<xml></xml>", ``)
		if err == nil || err.Error() != "SignXml: http error: Bad Request, status: 400" {
			t.Errorf("Expected error with unsuccessful status, got: %v", err)
		}
	})
}
