package models

type HttpRequest struct {
	Host       string
	Path       string
	Method     string
	Scheme     string
	Headers    map[string]string
	Parameters map[string]string
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{
		Headers:    make(map[string]string),
		Parameters: make(map[string]string),
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

func (r *HttpRequest) GetParameters() map[string]string {
	return r.Parameters
}
