package httpapi

import "errors"

var (
	ErrNotBearerToken = errors.New("not Bearer token")
	ErrInternalServer = errors.New("internal server error")
)
