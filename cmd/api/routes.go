package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Use(app.recoverPanic, app.rateLimit)
	r.Use(app.authenticate, app.enableCORS)

	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Get("/v1/update/makes", app.updateMakes)
	r.Get("/v1/update/models", app.updateModels)
	r.Get("/v1/update/generations", app.updateGenerations)
	r.Get("/v1/update/versions", app.updateVersions)
	r.Get("/v1/update/colors", app.updateColors)
	r.Get("/v1/update/transmissions", app.updateTransmissions)
	r.Get("/v1/update/body_types", app.updateBodyTypes)
	r.Get("/v1/update/fuel_types", app.updateFuelTypes)
	r.Get("/v1/update/cities", app.updateCities)
	r.Get("/v1/update/areas", app.updateAreas)
	r.Get("/v1/update/registrations", app.updateRegistration)

	r.Post("/v1/users/register", app.registerUserHandler)
	r.Get("/v1/users/getDetails", app.requireAuthenticatedUser(app.getUserHandler))
	r.Put("/v1/users/update", app.requireAuthenticatedUser(app.updateUserHandler))
	r.Post("/v1/users/updateProfilePic", app.requireAuthenticatedUser(app.updateProfilePicHandler))
	r.Post("/v1/users/authentication", app.createAuthenticationTokenHandler)

	r.Get("/v1/listings/{id}", app.getListingHandler)
	r.Post("/v1/listings", app.requireAuthenticatedUser(app.createListingHandler))
	r.Post("/v1/gallery", app.requireAuthenticatedUser(app.saveGalleryHandler))
	r.Get("/v1/listings", app.getListingsByFilter)
	r.Get("/v1/listings/home", app.getHomeFeed)

	r.Handle("/media/*", http.StripPrefix("/media/", http.FileServer(http.Dir("public/media"))))

	return r
}
