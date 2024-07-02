package models

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetHost(t *testing.T) {
	restful := NewHttpRequest(&RestfulConfiguration{
		CustomHost: "custom.example2.com",
	})
	restful.Host = "example.com"

	assert.Equal(t, restful.GetHost(), "custom.example2.com")
}
