package parsers

import (
	"github.com/tidwall/sjson"
	"github.com/whyayen/log2curl/internal/models"
	"slices"
	"strings"
)

type CloudWatchParser struct {
	Config *models.HttpRequestConfiguration
}

func NewCloudWatchParser(config *models.HttpRequestConfiguration) *CloudWatchParser {
	return &CloudWatchParser{
		Config: config,
	}
}

func (p *CloudWatchParser) Parse(log *map[string]string) *models.HttpRequest {
	restful := models.NewHttpRequest()
	parameters := ""

	for key, value := range *log {
		switch {
		case key == p.Config.Host:
			if p.Config.CustomHost != "" {
				restful.Host = p.Config.CustomHost
			} else {
				restful.Host = value
			}
		case key == p.Config.Path:
			restful.Path = value
		case key == p.Config.Method:
			restful.Method = value
		case key == p.Config.Scheme:
			restful.Scheme = value
		case strings.HasPrefix(key, p.Config.HeaderPrefix+"."):
			header := strings.TrimPrefix(key, p.Config.HeaderPrefix+".")
			if slices.Contains(p.Config.WhitelistHeaders, header) {
				restful.Headers[header] = value
			}
		case strings.HasPrefix(key, p.Config.ParameterPrefix+"."):
			newKey := strings.TrimPrefix(key, p.Config.ParameterPrefix+".")
			pJson, err := sjson.Set(parameters, newKey, value)

			if err != nil {
				restful.InvalidRequest = true
				return restful
			}
			parameters = pJson
		}
	}
	restful.Parameters = parameters

	return restful
}
