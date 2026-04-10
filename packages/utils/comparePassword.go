package utils

import "golang.org/x/crypto/bcrypt"

func Comparepassword(exist, new string) error{
	err := bcrypt.CompareHashAndPassword([]byte(exist), []byte(new))
	return err
}