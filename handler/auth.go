package handler

import (
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			usernameCookie, err1 := r.Cookie("username")
			tokenCookie, err2 := r.Cookie("token")
			if err1 != nil || err2 != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			username := usernameCookie.Value
			token := tokenCookie.Value

			if len(username) < 3 || !isTokenValid(username, token) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("fail"))
				return
			}
			h(w, r)
		})
}
