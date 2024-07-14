package handler

import (
	"ankasa-be/src"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()

	handler().ServeHTTP(w, r)
}

func handler() http.HandlerFunc {
	app := src.App()

	return adaptor.FiberApp(app)
}
