package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"backend/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	AuthService service.AuthServiceInterface
	Enforcer    *casbin.Enforcer
}

func NewAuthHandler(authService service.AuthServiceInterface, enforcer *casbin.Enforcer) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Enforcer:    enforcer,
	}
}

func (h *AuthHandler) OAuthLoginHandler(c *gin.Context) {
	var req dto.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Payload permintaan tidak valid",
			"detail":  err.Error(),
		})
		return
	}

	ipAddress := utils.GetIPAddress(c.Request)
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	_, err := h.AuthService.VerifyFirebaseIDToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Firebase tidak valid atau kadaluarsa",
			"detail":  err.Error(),
		})
		return
	}

	resp, err := h.AuthService.HandleOAuth(req.Token, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memproses login",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LoginManualHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Payload permintaan tidak valid",
			"detail":  err.Error(),
		})
		return
	}

	ipAddress := utils.GetIPAddress(c.Request)
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	resp, err := h.AuthService.LoginManual(req, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal memproses permintaan login",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyOTPHandler(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user id",
			"detail":  err.Error(),
		})
		return
	}

	ipAddress := utils.GetIPAddress(c.Request)
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request payload",
			"detail":  err.Error(),
		})
		return
	}

	resp, err := h.AuthService.VerifyOTP(userID, req.OTP, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized, invalid or expired OTP",
			"detail":  err.Error(),
		})
		return
	}

	log.Printf("[INFO] VerifyOTP | OTP berhasil diverifikasi \nuser id: %s \nIP: %s)", userID, ipAddress)
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) ResendOTPHandler(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user id tidak valid",
			"detail":  err.Error(),
		})
		return
	}

	ipAddress := utils.GetIPAddress(c.Request)
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	resp, err := h.AuthService.ResendOTP(userID, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengirim ulang OTP",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"detail":  "authorization header is missing",
		})
		return
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"detail":  "invalid authorization header format",
		})
		return
	}

	tokenStr := tokenParts[1]
	ipAddress := utils.GetIPAddress(c.Request)
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	if err := h.AuthService.Logout(tokenStr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "logout gagal",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout berhasil",
	})
}
