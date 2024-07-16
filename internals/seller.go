package data

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/iamgak/go-ecommerce/validator"
)

type SellerDB struct {
	db  *sql.DB
	ctx context.Context
}

func (s *SellerDB) CreateAccount(user *Seller) error {
	newHashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return ErrCantUseGeneratePassword
	}

	token := GenerateToken()
	var user_id int

	err = s.db.QueryRowContext(s.ctx, "INSERT INTO seller_listing (email, hashed_password,company_name, region_id, district_id, pincode, addr, pancard,mobile,activation_token) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id", user.Email, string(newHashedPassword), user.CompanyName, user.Region_id, user.District_id, user.Pincode, user.Addr, user.Pancard, user.Mobile, token).Scan(&user_id)
	if err == nil {
		err = s.SellerLog("Account Activate", user_id)
		return err
	}

	return err
}

func (s *SellerDB) ActivateAccount(token string) error {
	var user_id int

	err := s.db.QueryRowContext(s.ctx, "UPDATE seller_listing SET activation_token = NULL, active = TRUE WHERE activation_token = $1 RETURNING id", token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}

		return err
	}

	err = s.SellerLog("Account Activate", user_id)
	return err
}

func (s *SellerDB) ValidCredentials(user *User) (int, error) {
	var id int
	var hashedPassword string

	err := s.db.QueryRowContext(s.ctx, "SELECT id, hashed_password FROM seller_listing WHERE email = $1 AND active = TRUE", user.Email).Scan(&id, &hashedPassword)
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

func (s *SellerDB) Login(user *User) (string, error) {
	user_id, err := s.ValidCredentials(user)
	if err != nil {
		return "", err
	}

	token := GenerateToken()
	_, err = s.db.Exec("UPDATE seller_listing SET login_token = $1, last_login = CURRENT_TIMESTAMP WHERE id = $2", token, user_id)
	if err == nil {
		err = s.SellerLog("Seller Login", user_id)
	}

	return token, err
}

func (s *SellerDB) Logout(token string) error {
	var user_id int

	err := s.db.QueryRowContext(s.ctx, "UPDATE seller_listing SET login_token = NULL WHERE login_token = $1", token).Scan(&user_id)
	if err == nil {
		err = s.SellerLog("Seller LogOut", user_id)
	}
	return err
}

func (s *SellerDB) ValidToken(token string) (int, error) {
	var id int

	err := s.db.QueryRowContext(s.ctx, "SELECT id FROM seller_listing WHERE login_token = $1 AND active = TRUE", token).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return id, ErrNoRecord
	}

	return id, err
}

func (s *SellerDB) EmailExist(email string) (bool, error) {
	var validId int

	err := s.db.QueryRowContext(s.ctx, "SELECT 1 FROM seller_listing WHERE email = $1", email).Scan(&validId)
	return validId > 0, err
}

func (s *SellerDB) ResetPassword(email string) error {
	var user_id int

	err := s.db.QueryRowContext(s.ctx, "SELECT id FROM seller_listing WHERE email = $1 AND active = TRUE", email).Scan(&user_id)
	if err != nil {
		return err
	}

	uri := GenerateToken()
	_, _ = s.db.Exec("UPDATE user_forget_passw SET superseded = TRUE WHERE user_id = $1", user_id)
	_, err = s.db.Exec("INSERT INTO user_forget_passw (user_id,uri) VALUES ($1,$2)", user_id, uri)
	if err == nil {
		err = s.SellerLog("Forget Password Requested", user_id)
		return err
	}

	return err
}

func (s *SellerDB) ForgetPasswordErrorCheck(user *ForgetPassword) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.Email), "Email", "Email field Cannot be Empty")
	if user.Email != "" {
		validator.CheckField(validator.ValidEmail(user.Email), "Email", "Please, Enter Valid Email")
	}

	return validator.Errors
}

func (s *SellerDB) NewPasswordErrorCheck(user *NewPassword) map[string]string {
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

func (s *SellerDB) ResetPasswordURI(reset_token string) (int, error) {
	var user_id int

	err := s.db.QueryRowContext(s.ctx, "SELECT user_id FROM user_forget_passw WHERE uri = $1 AND superseded = FALSE", reset_token).Scan(&user_id)
	return user_id, err
}

func (s *SellerDB) NewPassword(user_id int, password string) error {
	_, err := s.db.Exec("UPDATE seller_listing SET hashed_password = $1 WHERE id = $2", password, user_id)
	if err == nil {
		err = s.SellerLog("Reset Password", user_id)
		return err
	}

	return err
}

func (s *SellerDB) SellerLog(activity string, uid int) error {
	_, err := s.db.Exec("UPDATE seller_log SET superseded = TRUE WHERE activity = $1 AND user_id = $2", activity, uid)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO seller_log ( activity,user_id) VALUES ( $1, $2)", activity, uid)
	return err
}

func (s *SellerDB) ErrorCheck(user *Seller) map[string]string {
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
