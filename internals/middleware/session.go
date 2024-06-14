package middleware

import (
	"gateway/internals/tools"
	md "gateway/model"
	"net/http"
)

var URLauth = "http://localhost:8081"

// for each request, check if the session is valid
func ValidSessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const action = "authorized"
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			w.WriteHeader(http.StatusMovedPermanently)
			w.Write([]byte("Redirect to login page"))
		}

		bodyData := md.RequestBody{
			Action: action,
			Body:   cookie.Value,
		}

		resp, err := tools.SendRequest(w, bodyData, http.MethodPost, URLauth)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode == http.StatusAccepted {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMovedPermanently)
			w.Write([]byte("Redirect to login page"))
		}

	})
}
