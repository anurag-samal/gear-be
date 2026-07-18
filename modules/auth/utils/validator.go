package auth

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate = validator.New()

func NormalizeRegisterRequest(req *RegisterRequest) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)
}

func NormalizeLoginRequest(req *LoginRequest) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
}

func ValidateRegisterRequest(req RegisterRequest) error {

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)
	if err := validate.Struct(req); err != nil {
		return err
	}

	if err := ValidatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func ValidateLoginRequest(req LoginRequest) error {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if err := validate.Struct(req); err != nil {
		return err
	}

	return nil
}
