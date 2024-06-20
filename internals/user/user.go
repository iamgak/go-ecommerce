package user

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/iamgak/go-ecommerce/validator"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDB struct {
	DB *sql.DB
}

func (u *UserDB) CreateUser(user *User) error {
	token := u.generateToken()
	_, err := u.DB.Exec("INSERT INTO users (email, password, activation_token) VALUES($1,$2,$3)", user.Email, user.Password, token)
	log.Printf("INSERT INTO users (email,password) VALUES(%s,%s)", user.Email, user.Password)
	return err
}

func (u *UserDB) LoginUser(user *User) error {
	token := u.generateToken()
	_, err := u.DB.Exec("INSERT INTO users (email, password, activation_token) VALUES($1,$2,$3)", user.Email, user.Password, token)
	log.Printf("INSERT INTO users (email,password) VALUES(%s,%s)", user.Email, user.Password)
	return err
}

func (u *UserDB) ActivateUser(token string) error {
	_, err := u.DB.Exec("UPDATE users SET activation_token = NULL, active = TRUE WHERE activation_token = $1", token)
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

func (u *UserDB) UserLogin(user *User) (string, error) {
	var id int
	err := u.DB.QueryRow("SELECT id FROM users WHERE email = $1 AND password = $2 AND active = TRUE", user.Email, user.Password).Scan(&id)
	if err != nil {
		return "", err
	}

	token := u.generateToken()
	_, err = u.DB.Exec("UPDATE users SET login_token = $1, last_login = CURRENT_TIMESTAMP WHERE id = $2", token, id)
	return token, err
}

func (u *UserDB) UserLogout(authHeader string) error {
	_, err := u.DB.Exec("UPDATE users SET login_token = NULL WHERE login_token = $1", authHeader)
	return err
}

func (u *UserDB) ValidToken(token string) (int, error) {
	var id int
	err := u.DB.QueryRow("SELECT id FROM users WHERE login_token = ?", token).Scan(&id)
	return id, err
}

func (u *UserDB) EmailExist(email string) (bool, error) {
	var id int
	err := u.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	return id > 0, err
}

func (u *UserDB) CreateLoginToken(uid int, r *http.Request) string {
	// later epoch ?? jwt
	token := fmt.Sprintf("%s-%d %d", r.RemoteAddr, uid, time.Second)
	return token
}

func (u *UserDB) generateToken() string {
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
