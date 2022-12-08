package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"test-exercise/api/constant"
	"test-exercise/api/dto"
	"test-exercise/api/rest/util"
)

var availableTypes = map[string]struct{}{
	"Corporations":        {},
	"NonProfit":           {},
	"Cooperative":         {},
	"Sole Proprietorship": {},
}

func ValidateCompany(f http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			util.ResponseInternalServerError(response)
			return
		}

		company := dto.Company{}
		if err := json.Unmarshal(body, &company); err != nil {
			util.ResponseBadRequest(response, "")
			return
		}

		if _, ok := availableTypes[company.Type]; !ok {
			util.ResponseBadRequest(response, "")
			return
		}

		ctx := context.WithValue(request.Context(), constant.CTX_COMPANY, company)
		f.ServeHTTP(response, request.WithContext(ctx))
	}) //handlerFunc
}
