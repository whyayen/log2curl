package generators

import (
	"github.com/magiconair/properties/assert"
	"github.com/whyayen/log2curl/internal/models"
	"testing"
)

func TestRequestUrl(t *testing.T) {
	httpRequest := &models.HttpRequest{
		Host:   "localhost",
		Method: "GET",
		Path:   "/users",
		Scheme: "https",
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Parameters:     "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}",
		InvalidRequest: false,
	}
	q := NewCurlGenerator(httpRequest)

	assert.Equal(t, q.RequestUrl(), "https://localhost/users?id=66838a1d337a8cdc830b439c&subscribed=true")
}

func TestRequestUrl_WhenMethodIsPost(t *testing.T) {
	httpRequest := &models.HttpRequest{
		Host:   "localhost",
		Method: "POST",
		Path:   "/users",
		Scheme: "https",
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Parameters:     "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}",
		InvalidRequest: false,
	}
	q := NewCurlGenerator(httpRequest)

	assert.Equal(t, q.RequestUrl(), "https://localhost/users")
}

func TestGenerate_WhenRequestIsInvalid(t *testing.T) {
	httpRequest := &models.HttpRequest{
		Host: "",
	}
	q := NewCurlGenerator(httpRequest)
	result, err := q.Generate()

	assert.Equal(t, result, "")
	assert.Equal(t, err.Error(), "invalid request")
}

func TestGenerate_WhenRequestIsGet(t *testing.T) {
	httpRequest := &models.HttpRequest{
		Host:   "localhost",
		Method: "GET",
		Path:   "/users",
		Scheme: "https",
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Parameters:     "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}",
		InvalidRequest: false,
	}
	q := NewCurlGenerator(httpRequest)
	result, _ := q.Generate()

	assert.Equal(t, result, "curl -X GET https://localhost/users?id=66838a1d337a8cdc830b439c&subscribed=true \\ \n -H 'Authorization: Bearer token'")
}

func TestGenerate_WhenRequestIsPost(t *testing.T) {
	httpRequest := &models.HttpRequest{
		Host:   "localhost",
		Method: "POST",
		Path:   "/users",
		Scheme: "https",
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Parameters:     "{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}",
		InvalidRequest: false,
	}
	q := NewCurlGenerator(httpRequest)
	result, _ := q.Generate()

	assert.Equal(t, result, "curl -X POST https://localhost/users \\ \n -H 'Authorization: Bearer token' \\ \n -d '{\"id\":\"66838a1d337a8cdc830b439c\",\"subscribed\":\"true\"}'")
}
