package controller

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"test-exercise/api/constant"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"

	"github.com/go-chi/chi/v5"
)

const (
	fmt_error_bad_id            = "id is not integer: %v"
	fmt_error_company_not_found = "company with id %v not found"
)

func GetCompany(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.ResponseBadRequest(response, fmt.Sprintf(fmt_error_bad_id, idStr))
		return
	}

	company, err := repository.Repo.GetCompany(id)
	if err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
		util.ResponseInternalServerError(response)
		return
	}

	if company == nil {
		util.ResponseNotFound(response, fmt.Sprintf(fmt_error_company_not_found, id))
		return
	}

	util.ResponseOK(response, company)
}
