package controller

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"test-exercise/api/constant"
	"test-exercise/api/dto"
	"test-exercise/api/mb"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func DeleteCompany(response http.ResponseWriter, request *http.Request) {
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
	} else if company == nil {
		util.ResponseNotFound(response, "")
		return
	}

	if err = repository.Repo.DeleteCompany(id); err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
		util.ResponseInternalServerError(response)
		return
	}

	user := request.Context().Value(constant.CTX_USER).(*dto.User)
	if err = mb.Kafka.Produce(viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT), &dto.Event{Method: http.MethodDelete, UserEmail: user.Email, CompanyName: company.Name}); err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	util.ResponseOK(response, nil)
}
