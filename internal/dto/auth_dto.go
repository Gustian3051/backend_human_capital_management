package dto

// ==========================
// REQUEST
// ==========================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OAuthLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required,len=6"`
}

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	CompanyName string `json:"companyName" binding:"required"`
}

// ==========================
// RESPONSE
// ==========================

type AuthResponse struct {
	AccessToken  string       `json:"accessToken,omitempty"`  // temp / full token
	RefreshToken string       `json:"refreshToken,omitempty"` // only full auth
	OtpRequest   bool         `json:"otpRequest"`
	NeedsProfile bool         `json:"needsProfile"`

	UserInfo    *UserInfo    `json:"userInfo,omitempty"`
	CompanyInfo *CompanyInfo `json:"companyInfo,omitempty"`
}

// ==========================
// USER INFO (SMART STRUCT)
// ==========================

type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`

	// 🔥 Optional (hanya setelah register)
	Name       string           `json:"name,omitempty"`
	EmployeeID string           `json:"employeeId,omitempty"`
	Role       string           `json:"role,omitempty"`
	Permission []PermissionInfo `json:"permission,omitempty"`
}

// ==========================
// COMPANY INFO
// ==========================

type CompanyInfo struct {
	ID            string `json:"id"`
	IsActive      bool   `json:"isActive"`
	IsTrial       bool   `json:"isTrial"`
	CompanyName   string `json:"companyName"`
	EmployeeCount int    `json:"employeeCount"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	PackageName   string `json:"packageName"`
}

// ==========================
// PERMISSION
// ==========================

type PermissionInfo struct {
	Name string `json:"name"`
}

// ==========================
// OTHER RESPONSE
// ==========================

type ResendOTPResponse struct {
	Message string `json:"message"`
}
