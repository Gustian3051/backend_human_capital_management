package jwt

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserID           string   `json:"user_id"`
	EmployeeID       string   `json:"employee_id"`
	CompanyID        string   `json:"company_id"`
	Role             string   `json:"role"`
	RoleKey          string   `json:"role_key"`
	Permissions      []string `json:"permissions"`
	
	jwt.RegisteredClaims
}

type ResetPasswordClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}