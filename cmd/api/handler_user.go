package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	token, err := app.User.UserLogin(input)
	if err != nil {
		if err == sql.ErrNoRows {
			app.Message(w, 200, "Invalid", "Incorrect Credentails")
			return
		}

		app.Message(w, 200, "Internal Error", "Internal Server Error")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	cookie := &http.Cookie{
		Name:    "ldata",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	app.Message(w, 200, "Message", "Login Successfull")
}

func (app *Application) UserRegister(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	inValid, err := app.User.EmailExist(input.Email)
	if err != nil && err != sql.ErrNoRows {
		app.InfoLog.Print(err)
		return
	}

	if inValid {
		app.Message(w, 200, "Email", "Email already Exist")
		return
	}

	err = app.User.CreateUser(input, false)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	app.sendJSONResponse(w, 200, input)
}

func (app *Application) SellerRegister(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	inValid, err := app.User.EmailExist(input.Email)
	if err != nil && err != sql.ErrNoRows {
		app.InfoLog.Print(err)
		return
	}

	if inValid {
		app.Message(w, 200, "Email", "Email already Exist")
		return
	}

	err = app.User.CreateUser(input, true)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	app.sendJSONResponse(w, 200, input)
}

func (app *Application) UserActivationToken(w http.ResponseWriter, r *http.Request) {
	activation_token := r.URL.Query().Get("activation_token")
	if r.Method != "GET" || activation_token == "" {
		http.NotFound(w, r)
		return
	}

	err := app.User.ActivateUser(activation_token)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Your Account has been Verified")
}

func (app *Application) UserForgetPassword(w http.ResponseWriter, r *http.Request) {
	var input *models.ForgetPassword
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ForgetPasswordErrorCheck(input)

	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.User.ResetPassword(input.Email)
	if err != nil {
		app.InfoLog.Print(err)
		app.Message(w, 500, "Server Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "If your Email is registered you will get mail on given email account")
}

func (app *Application) NewPassword(w http.ResponseWriter, r *http.Request) {
	reset_token := r.URL.Query().Get("reset_token")
	if reset_token == "" {
		http.NotFound(w, r)
		return
	}

	user_id, err := app.User.PasswordURI(reset_token)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var input *models.NewPassword
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.Message(w, 500, "Server Error", "Internal Server Error")
	}

	validator := app.User.NewPasswordErrorCheck(input)
	if len(validator) > 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.User.NewPassword(user_id, input.Password)
	if err != nil {
		app.Message(w, 500, "Server Error", "Internal Server Error")
	}

	app.sendJSONResponse(w, 200, "Password Reset Successfully")
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ldata")
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
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

	if cookie.Value != "" && len(cookie.Value) == 40 {
		err = app.User.UserLogout(cookie.Value)
		if err != nil {
			app.Message(w, 500, "Internal Error", "Internal Server Error")
			return
		}
	}

	app.Message(w, 500, "Message", "Logout Successfully")
}
