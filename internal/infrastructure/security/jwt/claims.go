package jwt

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserID           string   `json:"user_id"`
	EmployeeID       string   `json:"employee_id"`
	CompanyID        string   `json:"company_id"`
	RoleID           string   `json:"role_id"`
	Role             string   `json:"role"`
	Permissions      []string `json:"permissions"`
	jwt.RegisteredClaims
}

type ResetPasswordClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}