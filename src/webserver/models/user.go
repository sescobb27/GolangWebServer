package models

import (
	"os"
	"regexp"
)

type User struct {
	Username *string
	Email    *string
	Age      uintptr
	Image    *os.File
	Password *string
}

const (
	EMAIL_REGEX    = `(?i)\A[(A-Z0-9)+\.\_\%\+\-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}\z`
	NAME_REGEX     = `\A[a-zA-Z]{3,20}\z`
	USERNAME_REGEX = `\A[a-zA-Z0-9]{3,20}\z`
	PASS_REGEX     = `\A([a-zA-Z0-9]|\.|\:|\!|\@|\#|\$|\%|\^|\&|\/|\*){6,20}\z`
)

func validate(regex string, input interface{}) (bool, error) {
	var (
		match bool
		err   error
	)
	match, err = regexp.MatchString(regex, input.(string))
	if err != nil {
		return false, err
	} else if match {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *User) Save() error {
	return nil
}

func ValidateEmail(email string) (bool, error) {
	return validate(EMAIL_REGEX, email)
}

func ValidateUsername(username string) (bool, error) {
	return validate(NAME_REGEX, username)
}

func ValidatePassword(password string) (bool, error) {
	return validate(PASS_REGEX, password)
}
