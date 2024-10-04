package goncanode

import (
	"github.com/nbah1990/goncanode/types"
	"testing"
	"time"

	"github.com/nbah1990/goncanode/entities"
)

func TestCreate(t *testing.T) {
	t.Run("DefaultVersion", func(t *testing.T) {
		options := entities.Options{
			ServiceUrl: "https://example.com",
			P12base64:  "base64string",
			P12pass:    "password",
			Timeout:    5 * time.Second,
			Version:    nil,
		}

		handler := Create(options)

		if _, ok := handler.(*NCANodeV1Handler); !ok {
			t.Errorf("Expected handler type *NCANodeV1Handler, got %T", handler)
		}

		v1Handler := handler.(*NCANodeV1Handler)
		if v1Handler.P12base64 != options.P12base64 {
			t.Errorf("Expected P12base64 %s, got %s", options.P12base64, v1Handler.P12base64)
		}
		if v1Handler.P12pass != options.P12pass {
			t.Errorf("Expected P12pass %s, got %s", options.P12pass, v1Handler.P12pass)
		}
		if v1Handler.Timeout != options.Timeout {
			t.Errorf("Expected Timeout %v, got %v", options.Timeout, v1Handler.Timeout)
		}
	})

	t.Run("VersionV1", func(t *testing.T) {
		version := types.NCAnodeV10
		options := entities.Options{
			ServiceUrl: "https://example.com",
			P12base64:  "base64string",
			P12pass:    "password",
			Timeout:    5 * time.Second,
			Version:    &version,
		}

		handler := Create(options)

		if _, ok := handler.(*NCANodeV1Handler); !ok {
			t.Errorf("Expected handler type *NCANodeV1Handler, got %T", handler)
		}

		v1Handler := handler.(*NCANodeV1Handler)
		if v1Handler.P12base64 != options.P12base64 {
			t.Errorf("Expected P12base64 %s, got %s", options.P12base64, v1Handler.P12base64)
		}
		if v1Handler.P12pass != options.P12pass {
			t.Errorf("Expected P12pass %s, got %s", options.P12pass, v1Handler.P12pass)
		}
		if v1Handler.Timeout != options.Timeout {
			t.Errorf("Expected Timeout %v, got %v", options.Timeout, v1Handler.Timeout)
		}
	})

	t.Run("VersionV3", func(t *testing.T) {
		version := types.NCAnodeV30
		options := entities.Options{
			ServiceUrl: "https://example.com",
			P12base64:  "base64string",
			P12pass:    "password",
			Timeout:    5 * time.Second,
			Version:    &version,
		}

		handler := Create(options)

		if _, ok := handler.(*NCANodeV3Handler); !ok {
			t.Errorf("Expected handler type *NCANodeV3Handler, got %T", handler)
		}

		v3Handler := handler.(*NCANodeV3Handler)
		if v3Handler.P12base64 != options.P12base64 {
			t.Errorf("Expected P12base64 %s, got %s", options.P12base64, v3Handler.P12base64)
		}
		if v3Handler.P12pass != options.P12pass {
			t.Errorf("Expected P12pass %s, got %s", options.P12pass, v3Handler.P12pass)
		}
		if v3Handler.Timeout != options.Timeout {
			t.Errorf("Expected Timeout %v, got %v", options.Timeout, v3Handler.Timeout)
		}
	})

	t.Run("UnknownVersion", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic on unknown version, but no panic occurred")
			} else {
				err, ok := r.(error)
				if !ok {
					t.Errorf("Expected error type in panic, got %T", r)
				} else {
					if err.Error() != "unknown version" {
						t.Errorf("Expected error message 'unknown version', got '%s'", err.Error())
					}
				}
			}
		}()

		unknownVersion := types.Version("unknown")
		options := entities.Options{
			ServiceUrl: "https://example.com",
			P12base64:  "base64string",
			P12pass:    "password",
			Timeout:    5 * time.Second,
			Version:    &unknownVersion,
		}

		_ = Create(options)
	})
}
