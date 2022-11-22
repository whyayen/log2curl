package convertors

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

func Curl(payload *string) string {
	method := gjson.Get(*payload, "method").String()
	parameters := gjson.Get(*payload, "parameters")
	headers := gjson.Get(*payload, "headers")
	url := getUrl(payload)

	command := fmt.Sprintf("curl -X %s %s", method, url)

	if headers.Exists() {
		headers.ForEach(func(key, value gjson.Result) bool {
			command += fmt.Sprintf(" \\ \n -H '%s: %s'", key.String(), value.String())
			return true
		})
	}

	if strings.ToUpper(method) != "GET" && parameters.Exists() {
		command += fmt.Sprintf(" \\ \n -d '%s'", parameters.String())
	}

	return command
}

func getUrl(payload *string) string {
	scheme := gjson.Get(*payload, "scheme").String()
	host := gjson.Get(*payload, "host").String()
	path := gjson.Get(*payload, "path").String()
	method := gjson.Get(*payload, "method").String()
	parameters := gjson.Get(*payload, "parameters")

	requestUrl := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	if strings.ToUpper(method) == "GET" && parameters.Exists() {
		q := requestUrl.Query()

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
