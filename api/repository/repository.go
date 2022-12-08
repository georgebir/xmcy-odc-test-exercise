package repository

import (
	"test-exercise/api/dto"
)

var Repo Repository

type Repository interface {
	GetUser(string) (*dto.User, error)
	GetCompany(int) (*dto.Company, error)
	CreateCompany(*dto.Company) (*dto.Company, error)
	UpdateCompany(*dto.Company) error
	DeleteCompany(int) error
	AddEvent(event *dto.Event) error
}
