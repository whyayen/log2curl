package parsers

import (
	"github.com/whyayen/log2curl/internal/models"
	"slices"
	"strings"
)

type CloudWatchParser struct {
	Config *models.RestfulConfiguration
}

func NewCloudWatchParser(config *models.RestfulConfiguration) *CloudWatchParser {
	return &CloudWatchParser{
		Config: config,
	}
}

func (p *CloudWatchParser) Parse(log *map[string]string) *models.HttpRequest {
	restful := models.NewHttpRequest(p.Config)

	for key, value := range *log {
		switch {
		case key == p.Config.Host:
			restful.Host = value
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
			restful.Parameters[strings.TrimPrefix(key, p.Config.ParameterPrefix+".")] = value
		}
	}

	return restful
}
