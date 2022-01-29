package middlewares

import (
	"net/http"
	"strings"

	"github.com/ShauryaAg/ProductAPI/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			utils.Error(w, r, "Malformed Token", http.StatusUnauthorized)
		} else {
			jwtToken := authHeader[1]
			decoded, err := utils.ParseToken(jwtToken)
			if err != nil || decoded == nil {
				utils.Error(w, r, err.Error(), http.StatusUnauthorized)
			} else {
				r.Header.Set("decoded", (*decoded)["userId"].(string))
				next.ServeHTTP(w, r)
			}
		}
	})
}
