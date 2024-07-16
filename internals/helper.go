package data

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func CreateLoginToken(uid int, r *http.Request) string {
	token := fmt.Sprintf("%s-%d %d", r.RemoteAddr, uid, time.Second)
	return token
}

func GenerateToken() string {
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

func GeneratePassword(newPassword string) ([]byte, error) {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	fmt.Print(newHashedPassword)
	return newHashedPassword, err
}

func Matches(hash []byte, plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
