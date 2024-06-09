package types

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(nick string, email string, password string) (acc *Account) {

	hashedPassword := HashingPassword(password)
	acc = &Account{
		Nickname:  nick,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	//oblsuga unique Nickname, email
	return acc

}

func HashingPassword(password string) string {
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(pwd)
}
