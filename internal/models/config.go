package models

type HttpRequestConfiguration struct {
	ParameterPrefix  string
	HeaderPrefix     string
	Path             string
	Host             string
	Method           string
	Scheme           string
	WhitelistHeaders []string
	CustomHost       string
}
