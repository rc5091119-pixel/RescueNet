package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashpasss,err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "",err
	}
	return hashpasss,nil
}

func CheckPasswordHash(password,hashpassword string)(bool,error){
	match,err := argon2id.ComparePasswordAndHash(password,hashpassword)
	if err != nil {
		return false,err
	}
	return match,nil
}
