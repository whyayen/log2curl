package models

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetHost(t *testing.T) {
	request := NewHttpRequest()
	request.Host = "example.com"

	assert.Equal(t, request.GetHost(), "example.com")
}

func TestGetPath(t *testing.T) {
	request := NewHttpRequest()
	request.Path = "/users"

	assert.Equal(t, request.GetPath(), "/users")
}

func TestGetMethod(t *testing.T) {
	request := NewHttpRequest()
	request.Method = "GET"

	assert.Equal(t, request.GetMethod(), "GET")
}

func TestGetScheme(t *testing.T) {
	request := NewHttpRequest()
	request.Scheme = "https"

	assert.Equal(t, request.GetScheme(), "https")
}

func TestGetHeaders(t *testing.T) {
	request := NewHttpRequest()
	request.Headers = map[string]string{
		"Authorization": "Bearer token",
	}

	assert.Equal(t, request.GetHeaders(), map[string]string{
		"Authorization": "Bearer token",
	})
}

func TestGetParameters(t *testing.T) {
	request := NewHttpRequest()
	request.Parameters = "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}"

	assert.Equal(t, request.GetParameters(), "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}")
}

func TestIsValid_WhenHostIsNil(t *testing.T) {
	request := NewHttpRequest()
	request.Method = "GET"
	request.Scheme = "https"

	assert.Equal(t, request.IsValid(), false)
}

func TestIsValid_WhenHostIsEmpty(t *testing.T) {
	request := NewHttpRequest()
	request.Host = ""
	request.Method = "GET"
	request.Scheme = "https"

	assert.Equal(t, request.IsValid(), false)
}

func TestIsValid_WhenInvalidRequestIsTrue(t *testing.T) {
	request := NewHttpRequest()
	request.Host = "example.com"
	request.Method = "GET"
	request.Scheme = "https"
	request.InvalidRequest = true

	assert.Equal(t, request.IsValid(), false)

}

func TestIsValid_WhenHostAndMethodAndSchemeAreNotNil(t *testing.T) {
	request := NewHttpRequest()
	request.Host = "example.com"
	request.Method = "GET"
	request.Scheme = "https"

	assert.Equal(t, request.IsValid(), true)
}

func TestIsGet_WhenMethodIsGet(t *testing.T) {
	request := NewHttpRequest()
	request.Method = "GET"

	assert.Equal(t, request.IsGet(), true)
}

func TestIsGet_WhenMethodIsNotGet(t *testing.T) {
	request := NewHttpRequest()
	request.Method = "POST"

	assert.Equal(t, request.IsGet(), false)
}
