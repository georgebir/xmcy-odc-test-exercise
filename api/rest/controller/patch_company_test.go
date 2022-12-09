package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-exercise/api/dto"
	"test-exercise/api/messagebroker"
	mb_mock "test-exercise/api/messagebroker/mock"
	"test-exercise/api/repository"
	repository_mock "test-exercise/api/repository/mock"
	"test-exercise/api/rest/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPatchCompany_HappyPath(t *testing.T) {
	expectedCompany := dto.Company{
		Id:                43,
		Name:              "name",
		Description:       nil,
		AmountOfEmployees: 25,
		Registered:        true,
		Type:              "Corporations",
	}

	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(&dto.User{Id: 1, Email: "email"}, nil).Once()
	repoMock.On("UpdateCompany", &expectedCompany).Return(nil).Once()
	repository.Repo = repoMock

	kafkaMock := &mb_mock.MessageBroker{}
	kafkaMock.On("Produce", mock.Anything, mock.Anything).Return(nil).Once()
	messagebroker.MBroker = kafkaMock

	body, err := json.Marshal(expectedCompany)
	assert.Nil(t, err)
	request, err := http.NewRequest(http.MethodPatch, "/companies", bytes.NewReader(body))
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Group(func(r chi.Router) {
			r.Use(middleware.ValidateCompany)
			r.Patch("/companies", PatchCompany)
		})
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusOK, rWriter.Code)
	repoMock.AssertExpectations(t)
	kafkaMock.AssertExpectations(t)
}
