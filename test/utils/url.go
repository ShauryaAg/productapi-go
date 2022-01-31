package utils

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetUrlParamInContext(r *http.Request, param string, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(param, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
