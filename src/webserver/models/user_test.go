package models

import (
	"testing"
)

// type User User

var (
	match bool
	err   error
)

func shouldMatch(match bool, err error, t *testing.T) {
	if !match || err != nil {
		t.Error("Validation should match but it didn't", match, err)
	}
}

func shoudNotMatch(match bool, err error, t *testing.T) {
	if match || err != nil {
		t.Error("Validation should not match but it did", match, err)
	}
}

func each(values []string, validation func(val string)) {
	for _, value := range values {
		validation(value)
	}
}

func TestValidateUsername(t *testing.T) {
	match_values := []string{"Sescob", "SeScoB", "sescob", "SESCOB"}
	no_match_values := []string{"", "v", "verylongnamewithmanyletters",
		"with spaces"}
	should_pass := func(val string) {
		match, err = ValidateUsername(val)
		shouldMatch(match, err, t)
	}
	should_not_pass := func(val string) {
		match, err = ValidateUsername(val)
		shoudNotMatch(match, err, t)
	}
	each(match_values, should_pass)
	each(no_match_values, should_not_pass)
}

func TestValidatePassword(t *testing.T) {
	match_values := []string{"abcDEF", "ABCDEF", "abcdef", "abcDEF23",
		"12345678", "23abcDEF", ".:!@#$%^&/", ":abcDEF:", "$:abcDEF:$",
		"%a$b@cD^EF", "/12$%45#!&", "*abs*dft"}
	no_match_values := []string{"", "r4", "abcDE", "with space",
		"<html></html>", "?¡¿|\"=", "verylongPasswordwithmanyletters"}
	should_pass := func(val string) {
		match, err = ValidatePassword(val)
		shouldMatch(match, err, t)
	}
	should_not_pass := func(val string) {
		match, err = ValidatePassword(val)
		shoudNotMatch(match, err, t)
	}
	each(match_values, should_pass)
	each(no_match_values, should_not_pass)
}

func TestValidateEmail(t *testing.T) {
	match_values := []string{"user_user@user.com", "_user@user.com",
		"_user@user.com", "user.user@user.com", "user_user@user.com",
		"_user_user@user.com", "user@user.com.co", "user1223@user.com",
		"123user@user.com", "1u2ser1@user.com", "1234@user.com",
		"user@12345.com", "user@user123.com", "user@5678user.com"}
	no_match_values := []string{"@user.com", "user@.com", "user@com",
		"user@user", "user..@user", "user..user.@user", "", "_.&@user.com",
		"user", "user@user.com.", "@.com", "user user@user.com",
		"user user @user.com"}

	should_pass := func(val string) {
		match, err = ValidateEmail(val)
		shouldMatch(match, err, t)
	}
	should_not_pass := func(val string) {
		match, err = ValidateEmail(val)
		shoudNotMatch(match, err, t)
	}
	each(match_values, should_pass)
	each(no_match_values, should_not_pass)
}
