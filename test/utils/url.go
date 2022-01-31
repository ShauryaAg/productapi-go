package utils

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SetUrlParamInContext sets the url param in the current request context
// and returns a new request
func SetUrlParamInContext(r *http.Request, param string, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(param, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
