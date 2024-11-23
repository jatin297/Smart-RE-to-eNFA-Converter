package user

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:              name,
		Email:             email,
		EncryptedPassword: password,
		CreatedAt:         time.Now().UTC(),
	}
}

func (u *User) ValidPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}
