package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-exercise/api/dto"
	"test-exercise/api/mb"
	mb_mock "test-exercise/api/mb/mock"
	"test-exercise/api/repository"
	repository_mock "test-exercise/api/repository/mock"
	"test-exercise/api/rest/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCompany_HappyPath(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(&dto.User{Id: 1, Email: "email", Name: "name"}, nil).Once()
	repoMock.On("GetCompany", 1).Return(&dto.Company{}, nil).Once()
	repoMock.On("DeleteCompany", 1).Return(nil).Once()
	repository.Repo = repoMock

	kafkaMock := &mb_mock.KafkaService{}
	kafkaMock.On("Produce", mock.Anything, mock.Anything).Return(nil).Once()
	mb.Kafka = kafkaMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusOK, rWriter.Code)
	repoMock.AssertExpectations(t)
	kafkaMock.AssertExpectations(t)
}

func TestDeleteCompany_IfDeleteError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(&dto.User{Id: 1, Email: "email", Name: "name"}, nil).Once()
	repoMock.On("GetCompany", 1).Return(&dto.Company{}, nil).Once()
	repoMock.On("DeleteCompany", 1).Return(errors.New("error")).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusInternalServerError, rWriter.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteCompany_IfCompanyNotFound_ReturnsError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(&dto.User{Id: 1, Email: "email", Name: "name"}, nil).Once()
	repoMock.On("GetCompany", 1).Return(nil, nil).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusNotFound, rWriter.Code)
	assert.Equal(t, fmt.Sprintf(fmt_error_company_not_found, 1), rWriter.Body.String())

	repoMock.AssertExpectations(t)
}

func TestDeleteCompany_IfGetCompanyError_ReturnsError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(&dto.User{Id: 1, Email: "email", Name: "name"}, nil).Once()
	repoMock.On("GetCompany", 1).Return(nil, errors.New("error")).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusInternalServerError, rWriter.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteCompany_IfAuthError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(nil, errors.New("error")).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusInternalServerError, rWriter.Code)
	repoMock.AssertExpectations(t)
}

func TestDeleteCompany_IfNotAuthorized(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetUser", "1234").Return(nil, nil).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodDelete, "/companies/1", nil)
	assert.Nil(t, err)
	request.Header.Set("access_token", "1234")
	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", DeleteCompany)
	})
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusUnauthorized, rWriter.Code)
	repoMock.AssertExpectations(t)
}
