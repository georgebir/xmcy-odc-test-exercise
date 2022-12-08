package controller

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"test-exercise/api/constant"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"

	"github.com/go-chi/chi/v5"
)

func GetCompany(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		util.ResponseBadRequest(response, "")
		return
	}

	company, err := repository.Repo.GetCompany(id)
	if err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
		util.ResponseInternalServerError(response)
		return
	}

	if company == nil {
		util.ResponseNotFound(response, "")
		return
	}

	util.ResponseOK(response, company)
}
