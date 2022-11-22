package parsers

import (
	"regexp"
	"strings"

	"github.com/tidwall/sjson"
)

type ConfigMapperInterface struct {
	ParameterPrefix  string
	HeaderPrefix     string
	Path             string
	Host             string
	Method           string
	Scheme           string
	WhitelistHeaders []string
	CustomHost       string
}

func General(configMap *ConfigMapperInterface, log *map[string]string) string {
	var payload string

	headerRegex := regexp.MustCompile(`^` + regexp.QuoteMeta(configMap.HeaderPrefix) + `\.+`)
	parameterRegex := regexp.MustCompile(`^` + regexp.QuoteMeta(configMap.ParameterPrefix) + `\.+`)

	for key, value := range *log {
		if key == configMap.Host {
			var host string
			if configMap.CustomHost != "" {
				host = configMap.CustomHost
			} else {
				host = value
			}

			payload, _ = sjson.Set(payload, "host", host)
			continue
		}

		if key == configMap.Path {
			payload, _ = sjson.Set(payload, "path", value)
			continue
		}

		if key == configMap.Method {
			payload, _ = sjson.Set(payload, "method", value)
			continue
		}

		if key == configMap.Scheme {
			payload, _ = sjson.Set(payload, "scheme", value)
			continue
		}

		isHeader := headerRegex.MatchString(key)
		if isHeader {
			newKey := headerRegex.ReplaceAllString(key, "")

			if contains(configMap.WhitelistHeaders, newKey) {
				payload, _ = sjson.Set(payload, "headers."+newKey, value)
				continue
			}
		}

		isParameter := parameterRegex.MatchString(key)
		if isParameter {
			newKey := parameterRegex.ReplaceAllString(key, "")
			payload, _ = sjson.Set(payload, "parameters."+newKey, value)
		}
	}

	return payload
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}

	return false
}
