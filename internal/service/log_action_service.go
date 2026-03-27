package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/common"
	"backend/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
	logger "backend/pkg/log"
)

type LogServiceInterface interface {
	LogAction(ctx context.Context, userID, companyID uuid.UUID, action, ip string)
}

type LogService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogServiceInterface {
	return &LogService{repo: repo}
}

func (s *LogService) LogAction(ctx context.Context, userID, companyID uuid.UUID, action, ip string) {

	employeeName := "New User"

	name, err := s.repo.GetEmployeeNameByUserID(ctx, userID)
	if err == nil && name != "" {
		employeeName = name
	}

	var companyProfileID *uuid.UUID
	if companyID != uuid.Nil {
		companyProfileID = &companyID
	}

	logEntry := &common.LogActionModel{
		ID:        uuid.New(),
		UserID:    &userID,
		CompanyID: companyProfileID,
		Action:    fmt.Sprintf("%s by %s", action, employeeName),
		IPAddress: ip,
		DateTime:  time.Now(),
	}

	if err := s.repo.Create(ctx, logEntry); err != nil {
		logger.Log.Error("Failed to create log", zap.Error(err))
	}
}