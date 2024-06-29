package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	models "github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) SellerRegister(w http.ResponseWriter, r *http.Request) {
	var input *models.Seller
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.Seller.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	Valid, err := app.Seller.EmailExist(input.Email)
	if err != nil && err != sql.ErrNoRows {
		app.ServerError(w, err)
		return
	}

	if Valid {
		app.Message(w, 200, "Email", "Email already Exist")
		return
	}

	err = app.Seller.CreateAccount(input)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	app.sendJSONResponse(w, 200, input)
}

func (app *Application) SellerLogin(w http.ResponseWriter, r *http.Request) {
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

	token, err := app.Seller.Login(input)
	if err != nil {
		if err == sql.ErrNoRows {
			app.Message(w, 200, "Invalid", "Incorrect Credentails")
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
	message := SuccessResp{Success: "Login Successfull"}
	app.sendJSONResponse(w, 200, message)
}

func (app *Application) SellerActivationToken(w http.ResponseWriter, r *http.Request) {
	activation_token := r.URL.Query().Get("activation_token")
	if r.Method != "GET" || activation_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Activation log or Not Get method")
		return
	}

	err := app.Seller.ActivateAccount(activation_token)
	if err != nil {
		if err == models.ErrNoRecord {
			app.NotFound(w)
			return
		}

		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "Your Account has been Verified")
}

func (app *Application) SellerForgetPassword(w http.ResponseWriter, r *http.Request) {
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

	err = app.Seller.ResetPassword(input.Email)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "If your Email is registered you will get mail on given email account")
}

func (app *Application) SellerNewPassword(w http.ResponseWriter, r *http.Request) {
	reset_token := r.URL.Query().Get("reset_token")
	if reset_token == "" {
		app.NotFound(w)
		app.InfoLog.Print("Empty Reset Token")
		return
	}

	user_id, err := app.Seller.ResetPasswordURI(reset_token)
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

	app.sendJSONResponse(w, 200, "Password Reset Successfully")
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

	err = app.User.Logout(cookie.Value)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 500, "Message", "Logout Successfully")
}
func (app *Application) ValidSeller(r *http.Request) (int, error) {
	cookie, err := r.Cookie("sldata")
	if err != nil || cookie.Value == "" {
		if err == http.ErrNoCookie {
			return 0, models.ErrNoCookieFound
		}

		return 0, err
	}

	// Validate the token using the app.Seller.ValidToken method
	id, err := app.Seller.ValidToken(cookie.Value)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return id, models.ErrNoRecord
		}

		return 0, err
	}

	return id, nil
}
