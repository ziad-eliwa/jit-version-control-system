package utils

import (
	"errors"
)


var (
	ErrInvalidUsername = errors.New("Invalid Username")
	ErrInvalidEmailAddres = errors.New("Invalid Email Address")
	ErrUserNameAlreadyExists = errors.New("Username already Exists")
)