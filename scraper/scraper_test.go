package scraper

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_urlToFilename(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
	}{
		{
			rawURL:   "https://pkg.go.dev/net/url#URL",
			expected: "pkg.go.dev.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			u, err := url.Parse(tt.rawURL)
			require.NoError(t, err)
			actual := urlToFilename(u)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
