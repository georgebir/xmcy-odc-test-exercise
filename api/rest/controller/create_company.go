package controller

import (
	"log"
	"net/http"
	"runtime/debug"

	"test-exercise/api/constant"
	"test-exercise/api/dto"
	"test-exercise/api/messagebroker"
	"test-exercise/api/repository"
	"test-exercise/api/rest/util"

	"github.com/spf13/viper"
)

func CreateCompany(response http.ResponseWriter, request *http.Request) {
	company, err := repository.Repo.CreateCompany(request.Context().Value(constant.CTX_COMPANY).(*dto.Company))
	if err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
		util.ResponseInternalServerError(response)
	}

	user := request.Context().Value(constant.CTX_USER).(*dto.User)
	if err = messagebroker.MBroker.Produce(viper.GetString(constant.KAFKA_TOPIC_ADD_EVENT), &dto.Event{Method: http.MethodPost, UserEmail: user.Email, CompanyName: company.Name}); err != nil {
		log.Printf(constant.FMT_ERROR, err, string(debug.Stack()))
	}

	util.ResponseOK(response, company)
}
