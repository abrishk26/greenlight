package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/abrishk26/greenlight/internal/data"
	"github.com/abrishk26/greenlight/internal/validator"
)

func (a *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name: input.Name,
		Email: input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddErr("email", "a user with this email already exists")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w,r , err)
		}
		return
	}

	token, err := a.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	fn := func() {

		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID": user.ID,
		}

		err = a.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			a.logger.PrintError(err, nil)
			return
		}
	}

	a.background(fn)

	err = a.writeJSON(w, map[string]any{ "user": user }, nil, http.StatusAccepted)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
