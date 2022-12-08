package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"test-exercise/api/constant"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"
)

func Auth(f http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("access_token")
		user, err := repository.Repo.GetUser(token)
		if err != nil {
			log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
			util.ResponseInternalServerError(response)
			return
		} else if user == nil {
			util.ResponseUnauthorized(response)
			return
		}

		ctx := context.WithValue(request.Context(), constant.CTX_USER, user)
		f.ServeHTTP(response, request.WithContext(ctx))
	}) //handlerFunc
}
