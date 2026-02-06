package encrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("密码不能为空")
	}
	
	const cost = 12
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	
	return string(hashBytes), nil
}

func VerifyPassword(hashedPassword, password string) (bool, error) {
	if hashedPassword == "" || password == "" {
		return false, errors.New("密码不能为空")
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return false, nil
	default:
		return false, err
	}
}