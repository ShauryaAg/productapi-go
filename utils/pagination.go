package utils

import (
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Pagination returns pagination options for mongo queries
func Pagination(r *http.Request, findOptions *options.FindOptions) (*options.FindOptions, error) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "4"
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		return nil, err
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return nil, err
	}

	if page == 1 {
		findOptions.SetSkip(0)
		findOptions.SetLimit(limit)
		return findOptions, nil
	}

	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
	return findOptions, nil
}
