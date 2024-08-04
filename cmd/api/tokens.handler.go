package main

import (
	"errors"
	"net/http"
	"time"

	"ghostprotocols.pk/internal/data"
	"ghostprotocols.pk/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Identifier == "" {
		v.AddError("identifier", "must be provided")
	} else {
		if validator.Matches(input.Identifier, validator.EmailRX) {
			data.ValidateEmail(v, input.Identifier)
		} else if validator.Matches(input.Identifier, validator.PhoneRX) {
			data.ValidatePhone(v, input.Identifier)
		} else {
			v.AddError("identifier", "must be a valid email address or phone number")
		}
	}

	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user *data.User
	if validator.Matches(input.Identifier, validator.EmailRX) {
		user, err = app.models.Users.GetByEmail(input.Identifier)
	} else if validator.Matches(input.Identifier, validator.PhoneRX) {
		user, err = app.models.Users.GetByPhone(input.Identifier)
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
