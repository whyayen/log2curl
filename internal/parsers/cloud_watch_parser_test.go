package parsers

import (
	"github.com/magiconair/properties/assert"
	"github.com/whyayen/log2curl/internal/models"
	"testing"
)

func TestParse(t *testing.T) {
	cfg := &models.RestfulConfiguration{
		Host:             "host",
		Path:             "path",
		Method:           "method",
		Scheme:           "scheme",
		HeaderPrefix:     "header",
		ParameterPrefix:  "parameter",
		WhitelistHeaders: []string{"Authorization"},
		CustomHost:       "custom.example2.com",
	}
	parser := NewCloudWatchParser(cfg)
	log := map[string]string{
		"host":                 "example.com",
		"header.Authorization": "Bearer token",
		"header.Content-Type":  "application/json",
		"header.User-Agent":    "Webkit 1.0",
		"parameter.id":         "66838a1d337a8cdc830b439c",
		"parameter.subscribed": "true",
		"path":                 "/users",
		"method":               "GET",
		"scheme":               "https",
	}
	result := parser.Parse(&log)

	assert.Equal(t, result.GetHost(), "custom.example2.com")
	assert.Equal(t, result.GetPath(), "/users")
	assert.Equal(t, result.GetMethod(), "GET")
	assert.Equal(t, result.GetScheme(), "https")
	assert.Equal(t, result.GetHeaders(), map[string]string{
		"Authorization": "Bearer token",
	})
	assert.Equal(t, result.GetParameters(), map[string]string{
		"id":         "66838a1d337a8cdc830b439c",
		"subscribed": "true",
	})
}
