package models

import "strings"

type HttpRequest struct {
	Host           string
	Path           string
	Method         string
	Scheme         string
	Headers        map[string]string
	Parameters     string
	InvalidRequest bool // flag for invalid request when parsing
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{
		Headers:        make(map[string]string),
		Parameters:     "",
		InvalidRequest: false,
	}
}

func (r *HttpRequest) GetHost() string {
	return r.Host
}

func (r *HttpRequest) GetPath() string {
	return r.Path
}

func (r *HttpRequest) GetMethod() string {
	return r.Method
}

func (r *HttpRequest) GetScheme() string {
	return r.Scheme
}

func (r *HttpRequest) GetHeaders() map[string]string {
	return r.Headers
}

func (r *HttpRequest) GetParameters() string {
	return r.Parameters
}

func (r *HttpRequest) IsValid() bool {
	return r.Host != "" && r.Method != "" && r.Scheme != "" && !r.InvalidRequest
}

func (r *HttpRequest) IsGet() bool {
	return strings.ToUpper(r.Method) == "GET"
}
