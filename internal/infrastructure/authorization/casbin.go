package authorization

import (
	"backend/pkg/log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitCasbin(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		logger.Log.Fatal("Failed to initialize casbin adapter", zap.Error(err))
	}

	enforcer, err := casbin.NewEnforcer("config/casbin_model.conf", adapter)
	if err != nil {
		logger.Log.Fatal("Failed to initialize casbin enforcer", zap.Error(err))
	}

	if err := enforcer.LoadPolicy(); err != nil {
		logger.Log.Fatal("Failed to load casbin policy", zap.Error(err))
	}

	logger.Log.Info("Casbin initialized successfully")

	return enforcer
}
