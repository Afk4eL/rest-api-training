package middlewares

import (
	"clean-rest-arch/internal/server/utils/csrf"
	httperrors "clean-rest-arch/internal/server/utils/http/http.errors"
	"net/http"
)

// TODO:logger
func CSRFValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(csrf.CSRFHeader)
		if err != nil {
			httperrors.SetErrResponse(w, r, err)

			return
		}

		tokenFromCookie := cookie.Value
		if tokenFromCookie == "" {
			httperrors.SetErrResponse(w, r, httperrors.ErrNoToken)

			return
		}

		tokenFromHeader := r.Header.Get(csrf.CSRFHeader)

		if tokenFromCookie != tokenFromHeader {
			httperrors.SetErrResponse(w, r, httperrors.ErrNoToken)

			return
		}

		next.ServeHTTP(w, r)
	})

}
