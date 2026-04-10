package utils

import "golang.org/x/crypto/bcrypt"

func Hashing(data string) (string, error){
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hashed), err
}