package Middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"webapptrials/Classroom/Auth/Utils"
	"webapptrials/Classroom/Shared"
)

func Authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, err := Utils.ValidateToken(tokenString)
		if err != nil {
			Shared.SendApiResponse[any](w, http.StatusUnauthorized, "Unauthorized", nil, "Invalid or expired token")
			return
		}
		next(w, r, ps)
	}
}
