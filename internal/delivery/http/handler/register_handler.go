package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	RegisterService service.RegisterServiceInterface
}

func NewRegisterHandler(registerService service.RegisterServiceInterface) *RegisterHandler {
	return &RegisterHandler{
		RegisterService: registerService,
	}
}

func (h *RegisterHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest

	// 🔹 1. Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	ipAddress := c.ClientIP()

	// 🔹 2. Call service (claims sudah di ctx)
	resp, err := h.RegisterService.Register(
		c.Request.Context(),
		req,
		ipAddress,
	)

	if err != nil {

		switch err.Error() {
		case "unauthorized":
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return

		case "invalid token for register":
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
			return

		case "user already completed profile":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to register",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}