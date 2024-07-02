package generators

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whyayen/log2curl/internal/models"
	"net/url"
)

type CurlGenerator struct {
	HttpRequest *models.HttpRequest
}

func NewCurlGenerator(httpRequest *models.HttpRequest) *CurlGenerator {
	return &CurlGenerator{
		HttpRequest: httpRequest,
	}
}

func (g *CurlGenerator) Generate() (string, error) {
	if !g.HttpRequest.IsValid() {
		return "", fmt.Errorf("invalid request")
	}

	command := fmt.Sprintf("curl -X %s %s", g.HttpRequest.Method, g.RequestUrl())
	for key, value := range g.HttpRequest.GetHeaders() {
		command += fmt.Sprintf(" \\ \n -H '%s: %s'", key, value)
	}

	if !g.HttpRequest.IsGet() && g.HttpRequest.GetParameters() != "" {
		command += fmt.Sprintf(" \\ \n -d '%s'", g.HttpRequest.GetParameters())
	}

	return command, nil
}

func (g *CurlGenerator) RequestUrl() string {
	requestUrl := url.URL{
		Scheme: g.HttpRequest.GetScheme(),
		Host:   g.HttpRequest.GetHost(),
		Path:   g.HttpRequest.GetPath(),
	}

	if g.HttpRequest.IsGet() && g.HttpRequest.GetParameters() != "" {
		q := requestUrl.Query()
		parameters := gjson.Parse(g.HttpRequest.GetParameters())

		parameters.ForEach(func(key, value gjson.Result) bool {
			if value.IsArray() {
				value.ForEach(func(_, childValue gjson.Result) bool {
					q.Add(key.String()+"[]", childValue.String())
					return true
				})
			} else if value.IsObject() {
				value.ForEach(func(childKey, childValue gjson.Result) bool {
					q.Add(key.String()+"["+childKey.String()+"]", childKey.String())
					return true
				})
			} else {
				q.Set(key.String(), value.String())
			}
			return true
		})
		requestUrl.RawQuery = q.Encode()
	}

	return requestUrl.String()
}
