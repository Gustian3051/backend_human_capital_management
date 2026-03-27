package http

import (
	"backend/config"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/infrastructure/database"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"firebase.google.com/go/v4/auth"

	jwtinfra "backend/internal/infrastructure/security/jwt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg *config.Config, enforcer *casbin.Enforcer, db *gorm.DB, redisClient *redis.Client, jwtService *jwtinfra.Service, firebaseAuthClient *auth.Client) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Health Check (public)
	r.GET("/health", handler.HealthCheck)

	// Swagger (hanya dev)
	if cfg.App.Env != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// ===== REPOSITORY =====
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	logRepo := repository.NewLogRepository(db)

	seeder := database.NewRolePermissionSeeder(permRepo)

	// ===== SERVICES DEPENDENCY =====
	permissionService := service.NewRoleAndPermissionService(permRepo, seeder, enforcer)
	companyService := service.NewCompanyService(companyRepo, subscriptionRepo)
	employeeService := service.NewEmployeeService(employeeRepo, permissionService)
	logService := service.NewLogService(logRepo)

	// ===== AUTH SERVICE =====
	authService := service.NewAuthService(
		permissionService,
		companyService,
		employeeService,
		logService,
		jwtService,
		redisClient,
		firebaseAuthClient,
		userRepo,
		cfg.Auth.DefaultPassword,
	)

	registerService := service.NewRegisterService(
		db,
		jwtService,
		redisClient,
		userRepo,
		companyRepo,
		employeeRepo,
		logService,
		enforcer,
		seeder,
	)

	registerHandler := handler.NewRegisterHandler(registerService)

	// ===== HANDLER =====
	authHandler := handler.NewAuthHandler(
		authService,
		enforcer,
	)

	// ===== PUBLIC ROUTES =====
	public := r.Group("/api/v1/auth")
	{
		public.POST("/oauth", authHandler.OAuthLoginHandler)
		public.POST("/login", authHandler.LoginManualHandler)

		public.POST("/verify-otp/:user_id", authHandler.VerifyOTPHandler)
		public.POST("/resend-otp/:user_id", authHandler.ResendOTPHandler)
	}

	// ===== PROTECTED ROUTES =====
	api := r.Group("/api/v1")

	api.Use(
		middleware.JWTMiddleware(jwtService, redisClient),
	)

	api.POST("/register", registerHandler.RegisterHandler)

	api.POST("/logout", authHandler.LogoutHandler)
	secured := api.Group("")
	secured.Use(middleware.RBACMiddleware(enforcer))
	{

		// secured.GET("/employees", ...)
	}

	return r
}
