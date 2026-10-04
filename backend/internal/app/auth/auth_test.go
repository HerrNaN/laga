package auth

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()

	t.Run("accepts secure and local origins", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			origin string
			rpID   string
		}{
			{
				name:   "HTTPS origin",
				origin: "https://example.com",
				rpID:   "example.com",
			},
			{
				name:   "localhost development",
				origin: "http://localhost:8080",
				rpID:   "localhost",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				handler, err := New(nil, test.rpID, test.origin, slog.Default())

				require.NoError(t, err)
				assert.NotNil(t, handler)
			})
		}
	})

	t.Run("rejects invalid origins", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			origin string
		}{
			{
				name:   "HTTP origin",
				origin: "http://example.com",
			},
			{
				name:   "localhost lookalike",
				origin: "http://localhost.example.com",
			},
			{
				name:   "credentials in origin",
				origin: "https://user@example.com",
			},
			{
				name:   "path in origin",
				origin: "https://example.com/path",
			},
			{
				name:   "query in origin",
				origin: "https://example.com?query=value",
			},
			{
				name:   "fragment in origin",
				origin: "https://example.com#fragment",
			},
			{
				name:   "missing host",
				origin: "https://",
			},
			{
				name:   "missing origin",
				origin: "",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				handler, err := New(nil, "example.com", test.origin, slog.Default())

				assert.Error(t, err)
				assert.Nil(t, handler)
			})
		}
	})

	t.Run("requires relying party id", func(t *testing.T) {
		t.Parallel()

		_, err := New(nil, "", "https://example.com", slog.Default())

		assert.Error(t, err)
	})

	t.Run("cookie is secure", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			ttl  time.Duration
		}{
			{
				name: "session cookie",
				ttl:  sessionTTL,
			},
			{
				name: "expired cookie",
				ttl:  -time.Second,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				handler := &Handler{}
				recorder := httptest.NewRecorder()

				handler.cookie(recorder, sessionCookie, "token", "/", test.ttl)

				response := recorder.Result()
				t.Cleanup(func() {
					assert.NoError(t, response.Body.Close())
				})

				cookies := response.Cookies()
				require.Len(t, cookies, 1)

				cookie := cookies[0]
				assert.True(t, cookie.Secure)
				assert.True(t, cookie.HttpOnly)
				assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
				assert.Equal(t, int(test.ttl.Seconds()), cookie.MaxAge)
			})
		}
	})
}
