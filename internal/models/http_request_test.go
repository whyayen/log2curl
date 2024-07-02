package models

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetHost(t *testing.T) {
	restful := NewHttpRequest()
	restful.Host = "example.com"

	assert.Equal(t, restful.GetHost(), "example.com")
}

func TestGetPath(t *testing.T) {
	restful := NewHttpRequest()
	restful.Path = "/users"

	assert.Equal(t, restful.GetPath(), "/users")
}

func TestGetMethod(t *testing.T) {
	restful := NewHttpRequest()
	restful.Method = "GET"

	assert.Equal(t, restful.GetMethod(), "GET")
}

func TestGetScheme(t *testing.T) {
	restful := NewHttpRequest()
	restful.Scheme = "https"

	assert.Equal(t, restful.GetScheme(), "https")
}

func TestGetHeaders(t *testing.T) {
	restful := NewHttpRequest()
	restful.Headers = map[string]string{
		"Authorization": "Bearer token",
	}

	assert.Equal(t, restful.GetHeaders(), map[string]string{
		"Authorization": "Bearer token",
	})
}

func TestGetParameters(t *testing.T) {
	restful := NewHttpRequest()
	restful.Parameters = map[string]string{
		"id":         "66838a1d337a8cdc830b439c",
		"subscribed": "true",
	}

	assert.Equal(t, restful.GetParameters(), map[string]string{
		"id":         "66838a1d337a8cdc830b439c",
		"subscribed": "true",
	})
}
