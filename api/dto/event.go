package dto

type Event struct {
	Method      string `json:"method"`
	UserEmail   string `json:"user_email"`
	CompanyName string `json:"company_name"`
}
