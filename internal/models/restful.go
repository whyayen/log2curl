package models

type Restful struct {
	Host                 string
	Path                 string
	Method               string
	Scheme               string
	Headers              map[string]string
	Parameters           map[string]string
	RestfulConfiguration *RestfulConfiguration
}

func NewRestful(cfg *RestfulConfiguration) *Restful {
	return &Restful{
		Headers:              make(map[string]string),
		Parameters:           make(map[string]string),
		RestfulConfiguration: cfg,
	}
}

func (r *Restful) GetHost() string {
	if r.RestfulConfiguration != nil && r.RestfulConfiguration.CustomHost != "" {
		return r.RestfulConfiguration.CustomHost
	}

	return r.Host
}

func (r *Restful) GetPath() string {
	return r.Path
}

func (r *Restful) GetMethod() string {
	return r.Method
}

func (r *Restful) GetScheme() string {
	return r.Scheme
}

func (r *Restful) GetHeaders() map[string]string {
	return r.Headers
}

func (r *Restful) GetParameters() map[string]string {
	return r.Parameters
}
