package models

import (
	"database/sql"
	"regexp"

	"github.com/iamgak/go-ecommerce/validator"
)

type SellerDB struct {
	DB *sql.DB
}

type SellerModelInterface interface {
	CreateAccount(*Seller) error
	ResetPassword(string) error
	ResetPasswordURI(string) (int, error)
	NewPassword(int, string) error
	ValidToken(string) (int, error)
	Login(*User) (string, error)
	EmailExist(string) (bool, error)
	Logout(string) error
	ActivateAccount(string) error
	ErrorCheck(*Seller) map[string]string
	ForgetPasswordErrorCheck(*ForgetPassword) map[string]string
	NewPasswordErrorCheck(*NewPassword) map[string]string
	SellerLog(string, int) error
}

func (u *SellerDB) CreateAccount(user *Seller) error {
	newHashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return ErrCantUseGeneratePassword
	}

	token := GenerateToken()
	var user_id int
	err = u.DB.QueryRow("INSERT INTO seller_listing (email, hashed_password,company_name, region_id, district_id, pincode, addr, pancard,mobile,activation_token) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id", user.Email, string(newHashedPassword), user.CompanyName, user.Region_id, user.District_id, user.Pincode, user.Addr, user.Pancard, user.Mobile, token).Scan(&user_id)
	if err == nil {
		err = u.SellerLog("Account Activate", user_id)
		return err
	}

	return err
}

func (u *SellerDB) ActivateAccount(token string) error {
	var user_id int
	err := u.DB.QueryRow("UPDATE seller_listing SET activation_token = NULL, active = TRUE WHERE activation_token = $1 RETURNING id", token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}

		return err
	}

	err = u.SellerLog("Account Activate", user_id)
	return err
}

func (u *SellerDB) ValidCredentials(user *User) (int, error) {
	var id int
	var hashedPassword string

	err := u.DB.QueryRow("SELECT id, hashed_password FROM seller_listing WHERE email = $1 AND active = TRUE", user.Email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNoRecord
		}

		return 0, err
	}

	valid, err := Matches([]byte(hashedPassword), user.Password)
	if !valid {
		return 0, ErrNoRecord
	}

	return id, err
}

func (u *SellerDB) Login(user *User) (string, error) {
	user_id, err := u.ValidCredentials(user)
	if err != nil {
		return "", err
	}

	token := GenerateToken()
	_, err = u.DB.Exec("UPDATE seller_listing SET login_token = $1, last_login = CURRENT_TIMESTAMP WHERE id = $2", token, user_id)
	if err == nil {
		err = u.SellerLog("Seller Login", user_id)
	}

	return token, err
}

func (u *SellerDB) Logout(token string) error {
	var user_id int
	err := u.DB.QueryRow("UPDATE seller_listing SET login_token = NULL WHERE login_token = $1", token).Scan(&user_id)
	if err == nil {
		err = u.SellerLog("Seller LogOut", user_id)
	}
	return err
}

func (u *SellerDB) ValidToken(token string) (int, error) {
	var id int
	err := u.DB.QueryRow("SELECT id FROM seller_listing WHERE login_token = $1 AND active = TRUE", token).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return id, ErrNoRecord
	}

	return id, err
}

func (u *SellerDB) EmailExist(email string) (bool, error) {
	var validId int
	err := u.DB.QueryRow("SELECT 1 FROM seller_listing WHERE email = $1", email).Scan(&validId)
	return validId > 0, err
}

func (u *SellerDB) ResetPassword(email string) error {
	var user_id int
	err := u.DB.QueryRow("SELECT id FROM seller_listing WHERE email = $1 AND active = TRUE", email).Scan(&user_id)
	if err != nil {
		return err
	}

	uri := GenerateToken()
	_, _ = u.DB.Exec("UPDATE user_forget_passw SET superseded = TRUE WHERE user_id = $1", user_id)
	_, err = u.DB.Exec("INSERT INTO user_forget_passw (user_id,uri) VALUES ($1,$2)", user_id, uri)
	if err == nil {
		err = u.SellerLog("Forget Password Requested", user_id)
		return err
	}

	return err
}

func (u *SellerDB) ForgetPasswordErrorCheck(user *ForgetPassword) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.Email), "Email", "Email field Cannot be Empty")
	if user.Email != "" {
		validator.CheckField(validator.ValidEmail(user.Email), "Email", "Please, Enter Valid Email")
	}

	return validator.Errors
}

func (u *SellerDB) NewPasswordErrorCheck(user *NewPassword) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.RepeatPassword), "RepeatPassword", "Repeat Password field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Password), "Password", "Password field Cannot be Empty")
	if user.Password != "" {
		validator.CheckField(validator.MaxChars(user.Password, 20), "Password", "Please, Enter Password length less than 20")
	}

	if user.Password != "" && user.RepeatPassword != "" {
		validator.CheckField(user.Password == user.RepeatPassword, "RepeatPassword", "Please, Enter Same Password")
	}

	return validator.Errors
}

func (u *SellerDB) ResetPasswordURI(reset_token string) (int, error) {
	var user_id int
	err := u.DB.QueryRow("SELECT user_id FROM user_forget_passw WHERE uri = $1 AND superseded = FALSE", reset_token).Scan(&user_id)
	return user_id, err
}

func (u *SellerDB) NewPassword(user_id int, password string) error {
	_, err := u.DB.Exec("UPDATE seller_listing SET hashed_password = $1 WHERE id = $2", password, user_id)
	if err == nil {
		err = u.SellerLog("Reset Password", user_id)
		return err
	}

	return err
}

func (u *SellerDB) SellerLog(activity string, uid int) error {
	_, err := u.DB.Exec("UPDATE seller_log SET superseded = TRUE WHERE activity = $1 AND user_id = $2", activity, uid)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("INSERT INTO seller_log ( activity,user_id) VALUES ( $1, $2)", activity, uid)
	return err
}

func (u *SellerDB) ErrorCheck(user *Seller) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.Email), "Email", "Email field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Password), "Password", "Password field Cannot be Empty")
	if user.Email != "" {
		validator.CheckField(validator.ValidEmail(user.Email), "Email", "Please, Enter Valid Email")
	}

	if user.Password != "" {
		validator.CheckField(validator.MaxChars(user.Password, 20), "Password", "Please, Enter Password length less than 20")
	}

	validator.CheckField(validator.NotBlank(user.CompanyName), "company_name", "CompanyName field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Addr), "addr", "addr field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Pincode), "pincode", "Pincode field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Pancard), "pancard", "Pancard field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Mobile), "mobile", "Mobile field Cannot be Empty")
	validator.CheckField(user.Region_id > 0, "region", "region field Cannot be Empty or zero")
	validator.CheckField(user.District_id > 0, "district", "District_id field Cannot be Empty Or 0")
	if len(user.Pincode) > 0 {
		PincodePattern := regexp.MustCompile(`^[0-9]{5,8}$`)
		validator.CheckField(PincodePattern.MatchString(user.Pincode), "Pincode", "Pincode Should contain Numeric value only and length between 5 to 8")
	}

	return validator.Errors
}
