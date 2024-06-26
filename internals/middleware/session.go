package middleware

import (
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"net/http"
	"time"
)

// for each request, check if the session is valid
func ValidSessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("sessionID")
		if err != nil {
			w.WriteHeader(http.StatusMovedPermanently)
			w.Write([]byte("Redirect to login page"))
		}

		if !time.Now().After(cookie.Expires) {
			next.ServeHTTP(w, r)
			return
		}

		const action = "logout"

		bodyData := md.RequestBody{
			Action: action,
			Body:   cookie.Value,
		}

		resp, err := tools.SendRequest(w, bodyData, http.MethodPost, conf.URLauth)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode == http.StatusOK {
			// Créez un nouveau cookie avec le même nom mais avec une expiration passée
			cookie := http.Cookie{Name: "sessionID", Expires: time.Unix(0, 0), Path: "/"}
			http.SetCookie(w, &cookie)

			w.WriteHeader(http.StatusMovedPermanently)
			w.Write([]byte("Redirect to login page"))
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	})
}
