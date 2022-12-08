package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-exercise/api/dto"
	"test-exercise/api/repository"
	repository_mock "test-exercise/api/repository/mock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetCompany_HappyPath_ReturnsCompany(t *testing.T) {
	expectedCompany := dto.Company{
		Id:                1,
		Name:              "name",
		Description:       nil,
		AmountOfEmployees: 23,
		Registered:        true,
		Type:              "type",
	}

	repoMock := &repository_mock.Repository{}
	repoMock.On("GetCompany", 1).Return(&expectedCompany, nil).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodGet, "/companies/1", nil)
	assert.Nil(t, err)

	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/companies/{id}", GetCompany)
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusOK, rWriter.Code)
	b := rWriter.Body.Bytes()
	var company dto.Company
	err = json.Unmarshal(b, &company)
	assert.Nil(t, err)
	assert.Equal(t, expectedCompany, company)

	repoMock.AssertExpectations(t)
}

func TestGetCompany_WithBadId_ReturnsError(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/companies/id", nil)
	assert.Nil(t, err)

	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/companies/{id}", GetCompany)
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusBadRequest, rWriter.Code)
	assert.Equal(t, fmt.Sprintf(fmt_error_bad_id, "id"), rWriter.Body.String())
}

func TestGetCompany_NotFound_ReturnsError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetCompany", 1).Return(nil, nil).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodGet, "/companies/1", nil)
	assert.Nil(t, err)

	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/companies/{id}", GetCompany)
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusNotFound, rWriter.Code)
	assert.Equal(t, fmt.Sprintf(fmt_error_company_not_found, 1), rWriter.Body.String())

	repoMock.AssertExpectations(t)
}

func TestGetCompany_IfInternalServerError_ReturnsError(t *testing.T) {
	repoMock := &repository_mock.Repository{}
	repoMock.On("GetCompany", 1).Return(nil, errors.New("error")).Once()
	repository.Repo = repoMock

	request, err := http.NewRequest(http.MethodGet, "/companies/1", nil)
	assert.Nil(t, err)

	rWriter := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/companies/{id}", GetCompany)
	r.ServeHTTP(rWriter, request)

	assert.Equal(t, http.StatusInternalServerError, rWriter.Code)
	repoMock.AssertExpectations(t)
}
