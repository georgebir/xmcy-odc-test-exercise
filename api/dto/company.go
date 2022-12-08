package dto

type Company struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	AmountOfEmployees int     `json:"amount_of_employees"`
	Registered        bool    `json:"registered"`
	Type              string  `json:"type"`
}
