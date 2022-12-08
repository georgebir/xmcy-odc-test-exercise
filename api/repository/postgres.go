package repository

import (
	"database/sql"
	"fmt"
	"runtime/debug"

	"test-exercise/api/constant"
	"test-exercise/api/dto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	postgres_host     = "postgres.host"
	postgres_port     = "postgres.port"
	postgres_user     = "postgres.user"
	postgres_name     = "postgres.name"
	postgres_password = "postgres.password"

	fmt_params = "host=%v port=%v user=%v dbname=%v password=%v sslmode=disable"
)

type postgresRepository struct {
	db                                                                                                 *sqlx.DB
	getUserStmt, getCompanyStmt, createCompanyStmt, updateCompanyStmt, deleteCompanyStmt, addEventStmt *sql.Stmt
}

func NewPostgresRepository() (*postgresRepository, error) {
	host := viper.GetString(postgres_host)
	port := viper.GetInt(postgres_port)
	user := viper.GetString(postgres_user)
	name := viper.GetString(postgres_name)
	password := viper.GetString(postgres_password)

	params := fmt.Sprintf(fmt_params, host, port, user, name, password)

	db, err := sqlx.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	pgRepo := &postgresRepository{
		db: db,
	}
	if pgRepo.getUserStmt, err = db.Prepare("select id, email, name from users where token = $1"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	if pgRepo.createCompanyStmt, err = db.Prepare("insert into companies(	name, description, amount_of_employees, registered, type) values($1, $2, $3, $4, $5) returning id"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	if pgRepo.updateCompanyStmt, err = db.Prepare("update companies set name = $1, description = $2, amount_of_employees = $3, registered = $4, type = $5 where id = $6"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	if pgRepo.getCompanyStmt, err = db.Prepare("select name, description, amount_of_employees, registered, type from companies where id = $1"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	if pgRepo.deleteCompanyStmt, err = db.Prepare("delete from companies where id = $1"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}
	if pgRepo.addEventStmt, err = db.Prepare("insert into events(method, user_email, company_name) values($1, $2, $3)"); err != nil {
		db.Close()
		return nil, fmt.Errorf(constant.FMT_ERROR, err, string(debug.Stack()))
	}

	return pgRepo, nil
}

func (r *postgresRepository) GetUser(token string) (*dto.User, error) {
	rows, err := r.getUserStmt.Query(token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var id int
	var email, name string
	if err = rows.Scan(&id, &email, &name); err != nil {
		return nil, err
	}

	return &dto.User{Id: id, Email: email, Name: name}, nil
}

func (r *postgresRepository) GetCompany(id int) (*dto.Company, error) {
	rows, err := r.getCompanyStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var name, cType string
	var description interface{}
	var amountOfEmployees int
	var registered bool

	if err = rows.Scan(&name, &description, &amountOfEmployees, &registered, &cType); err != nil {
		return nil, err
	}

	var pDescription *string
	if description != nil {
		pDescription = new(string)
		*pDescription = description.(string)
	}

	return &dto.Company{
		Id:                id,
		Name:              name,
		Description:       pDescription,
		AmountOfEmployees: amountOfEmployees,
		Registered:        registered,
		Type:              cType,
	}, nil
}

func (r *postgresRepository) CreateCompany(c *dto.Company) (*dto.Company, error) {
	rows, err := r.createCompanyStmt.Query(c.Name, c.Description, c.AmountOfEmployees, c.Registered, c.Type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = rows.Scan(&c.Id); err != nil {
		return nil, err
	}

	return c, nil
}

func (r *postgresRepository) UpdateCompany(c *dto.Company) error {
	_, err := r.updateCompanyStmt.Exec(c.Name, c.Description, c.AmountOfEmployees, c.Registered, c.Type, c.Id)
	return err
}

func (r *postgresRepository) DeleteCompany(id int) error {
	_, err := r.deleteCompanyStmt.Exec(id)
	return err
}

func (r *postgresRepository) AddEvent(event *dto.Event) error {
	_, err := r.addEventStmt.Exec(event.Method, event.UserEmail, event.CompanyName)
	return err
}
