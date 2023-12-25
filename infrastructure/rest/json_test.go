//go:build unit

package rest

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		assertFn func(t *testing.T, b string)
	}{
		{
			name: "should marshal json correctly",
			input: `{
	"email": "test@localhost.me",
	"password": "password"
}`,
			assertFn: func(t *testing.T, b string) {
				u := struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}{}
				reader := strings.NewReader(b)
				assert.NoError(t, ReadJSON(io.NopCloser(reader), &u))

				assert.Equal(t, "test@localhost.me", u.Email)
				assert.Equal(t, "password", u.Password)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertFn(t, tt.input)
		})
	}
}
