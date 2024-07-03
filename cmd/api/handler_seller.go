package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	models "github.com/iamgak/go-ecommerce/internals"
	"net/http"
	"time"
)

func (app *Application) SellerRegister(w http.ResponseWriter, r *http.Request) {
	var input *models.Seller
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.ErrorMessage(w, http.StatusInternalServerError, "Invalid JSON Payload")
		return
	}

	validator := app.Model.Sellers.ErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, 200, validator)
		return
	}

	Valid, err := app.Model.Sellers.EmailExist(input.Email)
	if err != nil && err != sql.ErrNoRows {
		app.ServerError(w, err)
		return
	}

	if Valid {
		app.ErrorMessage(w, 200, "Email already Exist")
		return
	}

	err = app.Model.Sellers.CreateAccount(input)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	app.FinalMessage(w, 200, input)
}

func (app *Application) SellerLogin(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.ErrorMessage(w, http.StatusInternalServerError, "Invalid JSON Payload")
		return
	}

	validator := app.Model.Users.ErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, 200, validator)
		return
	}

	token, err := app.Model.Sellers.Login(input)
	if err != nil {
		if err == sql.ErrNoRows {
			app.ErrorMessage(w, http.StatusNotFound, "Incorrect Credentials")
			return
		}

		app.ServerError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "sldata",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	app.FinalMessage(w, http.StatusAccepted, "Login Succesfull")
}

func (app *Application) SellerActivationToken(w http.ResponseWriter, r *http.Request) {
	activation_token := r.URL.Query().Get("activation_token")
	if r.Method != "GET" || activation_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Activation log or Not Get method")
		return
	}

	err := app.Model.Sellers.ActivateAccount(activation_token)
	if err != nil {
		if err == models.ErrNoRecord {
			app.NotFound(w)
			app.ErrorLog.Print(err)
			return
		}

		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Your Account has been Verified")
}

func (app *Application) SellerForgetPassword(w http.ResponseWriter, r *http.Request) {
	var input *models.ForgetPassword
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.CustomError(w, err, "Invalid JSON Payload", http.StatusInternalServerError)
		return
	}

	validator := app.Model.Users.ForgetPasswordErrorCheck(input)

	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	err = app.Model.Sellers.ResetPassword(input.Email)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "If your Email is registered you will get mail on given email account")
}

func (app *Application) SellerNewPassword(w http.ResponseWriter, r *http.Request) {
	reset_token := r.URL.Query().Get("reset_token")
	if reset_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Reset Token")
		return
	}

	user_id, err := app.Model.Sellers.ResetPasswordURI(reset_token)
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	var input *models.NewPassword
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.ErrorMessage(w, http.StatusInternalServerError, "Invalid JSON Payload")
		return
	}

	validator := app.Model.Users.NewPasswordErrorCheck(input)
	if len(validator) > 0 {
		app.ErrorMessage(w, 200, validator)
		return
	}

	err = app.Model.Users.NewPassword(user_id, input.Password)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, 200, "Password Reset Successfully")
}

func (app *Application) SellerLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sldata")
	if err != nil || cookie.Value == "" || len(cookie.Value) != 40 {
		app.NotFound(w)
		app.InfoLog.Print("Invalid Logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "ldata",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		// HttpOnly: true,
	})

	err = app.Model.Users.Logout(cookie.Value)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, 500, "Logout Successfully")
}

func (app *Application) ValidSeller(r *http.Request) (int, error) {
	cookie, err := r.Cookie("sldata")
	if err != nil || cookie.Value == "" {
		if err == http.ErrNoCookie {
			return 0, models.ErrNoCookieFound
		}

		return 0, err
	}

	// Validate the token using the app.Model.Sellers.ValidToken method
	id, err := app.Model.Sellers.ValidToken(cookie.Value)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return id, models.ErrNoRecord
		}
	}

	return id, err
}
