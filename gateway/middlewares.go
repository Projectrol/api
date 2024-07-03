package main

import (
	"context"
	"net/http"

	"github.com/lehoangvuvt/projectrol/common"
)

func (app *application) AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := r.Cookie("access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		tokenStr := cookieValue.Value
		claims, err := common.ParseToken(tokenStr, "access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		userId := -1
		for k, v := range claims {
			if k == "sub" {
				userId = int(v.(float64))
			}
		}
		if userId == -1 {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		ctx := context.WithValue(r.Context(), common.ContextUserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
