package middleware

import (
	"net/http"

	"backend/pkg/log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RBACMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ===== Get data from JWT Middleware =====
		userID := c.GetString("user_id")
		companyID := c.GetString("company_id")

		if userID == "" {
			logger.Log.Warn("RBAC: missing user_id")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		if companyID == "" {
			logger.Log.Warn("RBAC: missing company_id")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// ===== Casbin parameters =====
		sub := userID
		dom := companyID
		obj := c.FullPath()
		act := c.Request.Method

		// ===== Enforce =====
		allowed, err := enforcer.Enforce(sub, dom, obj, act)
		if err != nil {
			logger.Log.Error("RBAC enforcement error",
				zap.Error(err),
				zap.String("sub", sub),
				zap.String("dom", dom),
				zap.String("obj", obj),
				zap.String("act", act),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "authorization error",
			})
			c.Abort()
			return
		}

		if !allowed {
			logger.Log.Warn("RBAC forbidden",
				zap.String("sub", sub),
				zap.String("dom", dom),
				zap.String("obj", obj),
				zap.String("act", act),
			)

			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			c.Abort()
			return
		}

		logger.Log.Info("RBAC allowed",
			zap.String("sub", sub),
			zap.String("dom", dom),
			zap.String("obj", obj),
			zap.String("act", act),
		)

		c.Next()
	}
}