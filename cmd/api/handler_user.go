package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	models "github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	token, err := app.User.Login(input)
	if err != nil {
		if err == models.ErrNoRecord {
			message := ErrResp{Error: "Incorrect Credentials"}
			app.sendJSONResponse(w, 200, message)
			return
		}

		app.ServerError(w, err)
		return
	}

	// w.Header().Set("Authorization", "Bearer "+token)
	cookie := &http.Cookie{
		Name:    "ldata",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	message := SuccessResp{Success: "Login Successfull"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) UserRegister(w http.ResponseWriter, r *http.Request) {
	var input *models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
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

	err = app.User.CreateAccount(input)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	message := SuccessResp{Success: "Account Registered"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) UserActivationToken(w http.ResponseWriter, r *http.Request) {
	activation_token := r.URL.Query().Get("activation_token")
	if activation_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Activation token")
		return
	}

	err := app.User.ActivateAccount(activation_token)
	if err != nil {
		if err == models.ErrNoRecord {
			app.NotFound(w)
			return
		}

		app.ServerError(w, err)
		return
	}

	message := SuccessResp{Success: "Account Verified"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) UserForgetPassword(w http.ResponseWriter, r *http.Request) {
	var input *models.ForgetPassword
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.User.ForgetPasswordErrorCheck(input)

	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.User.ResetPassword(input.Email)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	message := SuccessResp{Success: "If Your Email's registered you will get link to reset password"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) NewPassword(w http.ResponseWriter, r *http.Request) {
	reset_token := r.URL.Query().Get("reset_token")
	if reset_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Reset Token")
		return
	}

	user_id, err := app.User.ResetPasswordURI(reset_token)
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	var input *models.NewPassword
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.User.NewPasswordErrorCheck(input)
	if len(validator) > 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.User.NewPassword(user_id, input.Password)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	message := SuccessResp{Success: "Password Reset Successfully"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ldata")
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

	err = app.User.Logout(cookie.Value)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	message := SuccessResp{Success: "Logout Successfully"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) ValidUser(r *http.Request) (int, error) {
	cookie, err := r.Cookie("ldata")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, models.ErrNoCookieFound
		}

		return 0, err
	}

	// Validate the token using the app.Seller.ValidToken method
	id, err := app.User.ValidToken(cookie.Value)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return id, models.ErrNoRecord
		}

		return 0, err
	}

	return id, err
}
