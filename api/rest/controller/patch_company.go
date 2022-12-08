package controller

import (
	"log"
	"net/http"
	"runtime/debug"

	"test-exercise/api/constant"
	"test-exercise/api/dto"
	"test-exercise/api/mb"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"

	"github.com/spf13/viper"
)

func PatchCompany(response http.ResponseWriter, request *http.Request) {
	company := request.Context().Value(constant.CTX_COMPANY).(*dto.Company)

	err := repository.Repo.UpdateCompany(company)
	if err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
		util.ResponseInternalServerError(response)
		return
	}

	user := request.Context().Value(constant.CTX_USER).(*dto.User)
	if err = mb.Kafka.Produce(viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT), &dto.Event{Method: http.MethodDelete, UserEmail: user.Email, CompanyName: company.Name}); err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
	}

	util.ResponseOK(response, company)
}
