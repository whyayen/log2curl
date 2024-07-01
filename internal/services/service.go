package services

import "context"

type ServiceInterface interface {
	getQueryResults(ctx context.Context) (interface{}, error)
}
