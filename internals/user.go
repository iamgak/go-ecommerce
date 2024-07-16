package data

import (
	"context"
	"database/sql"

	"github.com/iamgak/go-ecommerce/validator"
)

type UserDB struct {
	db  *sql.DB
	ctx context.Context
}

func (u *UserDB) Close() {
	u.db.Close()
}

func (u *UserDB) CreateAccount(user *User) error {
	token := GenerateToken()
	newHashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return ErrCantUseGeneratePassword
	}

	var user_id int
	err = u.db.QueryRowContext(u.ctx, "INSERT INTO user_listing (email, hashed_password, activation_token) VALUES($1,$2,$3) RETURNING id", user.Email, string(newHashedPassword), token).Scan(&user_id)
	if err != nil {
		return ErrCantAddUser
	}

	err = u.ActivityLog("User Registered", user_id)
	return err
}

func (u *UserDB) ActivateAccount(token string) error {
	var user_id int
	err := u.db.QueryRowContext(u.ctx, "UPDATE user_listing SET activation_token = NULL, active = TRUE WHERE activation_token = $1 RETURNING id", token).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}

		return err
	}

	err = u.ActivityLog("User Activated", user_id)
	return err
}

func (u *UserDB) ErrorCheck(user *User) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.Email), "Email", "Email field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user.Password), "Password", "Password field Cannot be Empty")
	if user.Email != "" {
		validator.CheckField(validator.ValidEmail(user.Email), "Email", "Please, Enter Valid Email")
	}

	if user.Password != "" {
		validator.CheckField(validator.MaxChars(user.Password, 20), "Password", "Please, Enter Password length less than 20")
	}

	return validator.Errors
}

func (u *UserDB) ValidCredentials(user *User) (int, error) {
	var id int
	var password string
	err := u.db.QueryRowContext(u.ctx, "SELECT id, hashed_password FROM user_listing WHERE email = $1 AND active = TRUE", user.Email).Scan(&id, &password)
	if err != nil && err == sql.ErrNoRows {
		return id, ErrNoRecord
	}

	valid, err := Matches([]byte(password), user.Password)
	if !valid {
		return 0, ErrNoRecord
	}

	return id, err
}

func (u *UserDB) Login(user *User) (string, error) {
	user_id, err := u.ValidCredentials(user)
	if err != nil {
		return "", err
	}

	token := GenerateToken()
	_, err = u.db.Exec("UPDATE user_listing SET login_token = $1, last_login = CURRENT_TIMESTAMP WHERE id = $2", token, user_id)
	if err == nil {
		_ = u.ActivityLog("User Login", user_id)
	}

	return token, err
}

func (u *UserDB) Logout(token string) error {
	var user_id int
	err := u.db.QueryRowContext(u.ctx, "UPDATE user_listing SET login_token = NULL WHERE login_token = $1 RETURNING id", token).Scan(&user_id)
	if err == nil {
		_ = u.ActivityLog("User Logout", user_id)
	}

	return err
}

func (u *UserDB) ValidToken(token string) (int, error) {
	var id int
	err := u.db.QueryRowContext(u.ctx, "SELECT id FROM user_listing WHERE login_token = $1 AND active = TRUE", token).Scan(&id)
	return id, err
}

func (u *UserDB) EmailExist(email string) (bool, error) {
	var id int
	err := u.db.QueryRowContext(u.ctx, "SELECT id FROM user_listing WHERE email = $1", email).Scan(&id)
	return id > 0, err
}

func (u *UserDB) ResetPassword(email string) error {
	var user_id int
	err := u.db.QueryRowContext(u.ctx, "SELECT id FROM user_listing WHERE email = $1", email).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}

		return err
	}

	uri := GenerateToken()
	_, err = u.db.Exec("UPDATE user_forget_passw SET superseded = TRUE WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}

	_, err = u.db.Exec("INSERT INTO user_forget_passw (user_id,uri) VALUES ($1,$2)", user_id, uri)
	if err == nil {
		err = u.ActivityLog("Reset Password Requested", user_id)
	}

	return err
}

func (u *UserDB) ForgetPasswordErrorCheck(user *ForgetPassword) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user.Email), "Email", "Email field Cannot be Empty")
	if user.Email != "" {
		validator.CheckField(validator.ValidEmail(user.Email), "Email", "Please, Enter Valid Email")
	}

	return validator.Errors
}

func (u *UserDB) NewPasswordErrorCheck(user *NewPassword) map[string]string {
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

func (u *UserDB) ResetPasswordURI(reset_token string) (int, error) {
	var user_id int
	err := u.db.QueryRowContext(u.ctx, "SELECT user_id FROM user_forget_passw WHERE uri = $1 AND superseded = FALSE", reset_token).Scan(&user_id)
	if err != nil && err == sql.ErrNoRows {
		return user_id, ErrNoRecord
	}

	return user_id, err
}

func (u *UserDB) NewPassword(user_id int, password string) error {
	_, err := u.db.Exec("UPDATE user_listing SET hashed_password = $1 WHERE id = $2", password, user_id)
	if err != nil {
		return err
	}

	_, err = u.db.Exec("UPDATE user_forget_passw SET superseded = TRUE WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}

	err = u.ActivityLog("Password Reset", user_id)
	return err
}

func (u *UserDB) ActivityLog(activity string, uid int) error {
	_, err := u.db.Exec("UPDATE user_log SET superseded = TRUE WHERE activity = $1 AND user_id = $2", activity, uid)
	if err != nil {
		return err
	}

	_, err = u.db.Exec("INSERT INTO user_log ( activity,user_id) VALUES ( $1, $2)", activity, uid)
	return err
}
func (pr *UserDB) AddrErrorCheck(user_addr *UserAddr) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(user_addr.Addr), "Addr", "Addr field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user_addr.Pincode), "Pincode", "Pincode field Cannot be Empty")
	validator.CheckField(validator.NotBlank(user_addr.Mobile), "Mobile", "Mobile field Cannot be Empty")
	validator.CheckField(user_addr.Region_id > 0, "Region_id", "Region_id field Cannot be Empty or zero")
	validator.CheckField(user_addr.District_id > 0, "District_id", "District_id field Cannot be Empty Or 0")
	if len(user_addr.Pincode) > 0 {
		validator.CheckField(validator.MatchString("^[0-9]{5,8}$", user_addr.Pincode), "Pincode", "Pincode Should contain Numeric value only")
	}

	return validator.Errors
}

func (u *UserDB) CreateAddr(addr *UserAddr, uid int) (int, error) {
	var AddrId int
	err := u.db.QueryRowContext(u.ctx, "INSERT INTO user_addr (user_id, mobile, region_id, district_id, pincode, addr) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id ", uid, addr.Mobile, addr.Region_id, addr.District_id, addr.Pincode, addr.Addr).Scan(&AddrId)
	return AddrId, err
}
