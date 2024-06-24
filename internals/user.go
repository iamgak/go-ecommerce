package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/iamgak/go-ecommerce/validator"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPassword struct {
	Email string `json:"email"`
}

type NewPassword struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UserDB struct {
	DB *sql.DB
}

type UserModelInterface interface {
	CreateUser(*User, bool) error
	ResetPassword(string) error
	PasswordURI(string) (int, error)
	NewPassword(int, string) error
	ValidToken(string) (int, bool)
	UserLogin(*User) (string, error)
	EmailExist(string) (bool, error)
	UserLogout(string) error
	ActivateUser(string) error
	ErrorCheck(*User) map[string]string
	ForgetPasswordErrorCheck(*ForgetPassword) map[string]string
	NewPasswordErrorCheck(*NewPassword) map[string]string
	GenerateToken() string
}

func (u *UserDB) CreateUser(user *User, seller bool) error {
	token := u.GenerateToken()
	newHashedPassword, err := u.GeneratePassword(user.Password)
	if err != nil {
		return ErrCantUseGeneratePassword
	}

	_, err = u.DB.Exec("INSERT INTO user_listing (email, password, activation_token,seller) VALUES($1,$2,$3,$4)", user.Email, string(newHashedPassword), token, seller)
	return err
}

func (u *UserDB) ActivateUser(token string) error {
	_, err := u.DB.Exec("UPDATE user_listing SET activation_token = NULL, active = TRUE WHERE activation_token = $1", token)
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
	newHashedPassword, err := u.GeneratePassword(user.Password)
	if err != nil {
		return id, err
	}
	err = u.DB.QueryRow("SELECT id FROM user_listing WHERE email = $1 AND password = $2 AND active = TRUE", user.Email, string(newHashedPassword)).Scan(&id)
	return id, err
}

func (u *UserDB) UserLogin(user *User) (string, error) {
	id, err := u.ValidCredentials(user)
	if err != nil {
		return "", err
	}

	token := u.GenerateToken()
	_, err = u.DB.Exec("UPDATE user_listing SET login_token = $1, last_login = CURRENT_TIMESTAMP WHERE id = $2", token, id)
	return token, err
}

func (u *UserDB) UserLogout(authHeader string) error {
	_, err := u.DB.Exec("UPDATE user_listing SET login_token = NULL WHERE login_token = $1", authHeader)
	return err
}

func (u *UserDB) ValidToken(token string) (int, bool) {
	var id int
	var seller bool
	_ = u.DB.QueryRow("SELECT id, seller FROM user_listing WHERE login_token = $1", token).Scan(&id, &seller)
	return id, seller
}

func (u *UserDB) EmailExist(email string) (bool, error) {
	var id int
	err := u.DB.QueryRow("SELECT id FROM user_listing WHERE email = $1", email).Scan(&id)
	return id > 0, err
}

func (u *UserDB) ResetPassword(email string) error {
	var id int
	err := u.DB.QueryRow("SELECT id FROM user_listing WHERE email = $1", email).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if id > 0 {
		uri := u.GenerateToken()
		_, _ = u.DB.Exec("UPDATE user_forget_passw SET superseded = TRUE WHERE user_id = $1", id)
		_, err = u.DB.Exec("INSERT INTO user_forget_passw (user_id,uri) VALUES ($1,$2)", id, uri)
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

func (u *UserDB) CreateLoginToken(uid int, r *http.Request) string {
	token := fmt.Sprintf("%s-%d %d", r.RemoteAddr, uid, time.Second)
	return token
}

func (u *UserDB) GenerateToken() string {
	// Get current epoch time (Unix timestamp in seconds)
	epoch := time.Now().Unix()

	// Generate a random number (for simplicity, let's use a random integer)
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000000) // Random number between 0 and 999999

	// Concatenate epoch time and random number
	tokenString := strconv.FormatInt(epoch, 10) + strconv.Itoa(randomNumber)

	// Hash the concatenated string using SHA-1
	hash := sha1.New()
	hash.Write([]byte(tokenString))
	hashed := hash.Sum(nil)

	// Convert hashed bytes to a hexadecimal string
	token := fmt.Sprintf("%x", hashed)

	return token
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

func (u *UserDB) PasswordURI(reset_token string) (int, error) {
	var user_id int
	err := u.DB.QueryRow("SELECT user_id FROM user_forget_passw WHERE uri = $1 AND superseded = FALSE", reset_token).Scan(&user_id)
	return user_id, err
}

func (u *UserDB) NewPassword(user_id int, password string) error {
	_, err := u.DB.Exec("UPDATE user_listing SET password = $1 WHERE id = $2", password, user_id)
	return err
}

func (u *UserDB) GeneratePassword(newPassword string) ([]byte, error) {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	return newHashedPassword, err
}

func (u *UserDB) ActivityLog(activity string, uid int64) {
	_, _ = u.DB.Exec("UPDATE `user_log` SET superseded = 1 WHERE activity = ? AND uid = ?", activity, uid)
	_, _ = u.DB.Exec("INSERT INTO `user_log` SET  activity = ? , uid = ?, superseded = 0", activity, uid)
}
