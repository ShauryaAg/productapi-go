package utils

import (
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func Pagination(r *http.Request, findOptions *options.FindOptions) *options.FindOptions {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "4"
	}

	page, _ := strconv.ParseInt(pageStr, 10, 32)
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	if page == 1 {
		findOptions.SetSkip(0)
		findOptions.SetLimit(limit)
		return findOptions
	}

	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
	return findOptions
}
