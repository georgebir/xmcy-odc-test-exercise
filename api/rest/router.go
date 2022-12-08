package rest

import (
	"fmt"
	"net/http"

	"test-exercise/api/rest/controller"
	"test-exercise/api/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func ListenAndServe(port int) {
	r := chi.NewRouter()

	r.Get("/companies/{id}", controller.GetCompany)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Delete("/companies/{id}", controller.DeleteCompany)

		r.Group(func(r chi.Router) {
			r.Use(middleware.ValidateCompany)
			r.Post("/companies", controller.CreateCompany)
			r.Patch("/companies", controller.PatchCompany)
		}) //group validate
	}) //group auth

	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
