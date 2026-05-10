package utils

import (
	"testing"
)

type TestUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"email"`
	Phone    string `json:"phone" binding:"phone"`
}

func TestValidateRequest_ValidData(t *testing.T) {
	RegisterCustomValidators()

	user := TestUser{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		Phone:    "13800138000",
	}

	errors := ValidateRequest(user)
	if len(errors) > 0 {
		t.Errorf("Expected no validation errors, got: %v", errors)
	}
}

func TestValidateRequest_MissingFields(t *testing.T) {
	RegisterCustomValidators()

	user := TestUser{}

	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected validation errors for missing fields")
	}

	// 检查是否包含必需字段错误
	found := false
	for _, err := range errors {
		if err.Field == "Username" || err.Field == "Password" {
			found = true
		}
	}
	if !found {
		t.Error("Expected errors for Username and Password fields")
	}
}

func TestValidateRequest_InvalidEmail(t *testing.T) {
	RegisterCustomValidators()

	user := TestUser{
		Username: "testuser",
		Password: "password123",
		Email:    "invalid-email",
	}

	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid email")
	}
}

func TestValidateRequest_InvalidPhone(t *testing.T) {
	RegisterCustomValidators()

	user := TestUser{
		Username: "testuser",
		Password: "password123",
		Phone:    "12345", // 太短，不是有效手机号
	}

	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid phone")
	}
}

func TestValidatePasswordStrong_Valid(t *testing.T) {
	type PwdUser struct {
		Password string `json:"password" binding:"password_strong"`
	}

	RegisterCustomValidators()

	// 有效密码：至少8位，包含字母和数字
	user := PwdUser{Password: "Pass12345"}
	errors := ValidateRequest(user)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid password, got: %v", errors)
	}
}

func TestValidatePasswordStrong_TooShort(t *testing.T) {
	type PwdUser struct {
		Password string `json:"password" binding:"password_strong"`
	}

	RegisterCustomValidators()

	user := PwdUser{Password: "Pass1"}
	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected error for password too short")
	}
}

func TestValidatePasswordStrong_NoNumber(t *testing.T) {
	type PwdUser struct {
		Password string `json:"password" binding:"password_strong"`
	}

	RegisterCustomValidators()

	user := PwdUser{Password: "Password"}
	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected error for password without number")
	}
}

func TestValidatePasswordStrong_NoLetter(t *testing.T) {
	type PwdUser struct {
		Password string `json:"password" binding:"password_strong"`
	}

	RegisterCustomValidators()

	user := PwdUser{Password: "12345678"}
	errors := ValidateRequest(user)
	if len(errors) == 0 {
		t.Error("Expected error for password without letter")
	}
}
