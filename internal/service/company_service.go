package service

import (
	"backend/internal/domain/company"
	"backend/internal/domain/subscription"
	"backend/internal/domain/user"
	"backend/internal/dto"
	"backend/internal/repository"
	"backend/pkg/helpers"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GeneralInformation(ctx context.Context, userID, companyID uuid.UUID, req *dto.GeneralInformationUpdate, ipAddress string) (*dto.CompanyProfileResponse, error)
// CompanyInformation(ctx context.Context, userID, companyID uuid.UUID, req *dto.CompanyInformationUpdate, ipAddress string) (*dto.CompanyProfileResponse, error)
// OwnerInformation(ctx context.Context, userID, companyID uuid.UUID, req *dto.OwnerInformationUpdate, ipAddress string) (*dto.CompanyProfileResponse, error)
// UploadCompanyDocument(ctx context.Context, userID, companyID uuid.UUID, req *dto.FileDocument, ipAddress string) (*dto.CompanyProfileResponse, error)
// DeleteCompanyDocument(ctx context.Context, userID, companyID, fileID uuid.UUID, ipAddress string) (*dto.CompanyProfileResponse, error)

// DataFilesCompany(ctx context.Context, userID, companyID uuid.UUID) ([]dto.FileDocumentData, error)
// DataCompanyProfile(ctx context.Context, companyID uuid.UUID) (*dto.CompanyProfileResponse, error)

type CompanyServiceInterface interface {
	GetCompanyForUser(ctx context.Context, user *user.UserModel) (*dto.CompanyInfo, error)
	CreateDefaultCompany(ctx context.Context, name string) (*company.CompanyModel, error)
}

type CompanyService struct {
	CompanyRepo repository.CompanyRepository
	SubscriptionRepo repository.SubscriptionRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository, subscriptionRepo repository.SubscriptionRepository) CompanyServiceInterface {
	return &CompanyService{
		CompanyRepo: companyRepo,
		SubscriptionRepo: subscriptionRepo,
	}
}

func (s *CompanyService) generateCompanyCode(tx *gorm.DB) (string, error) {
    const prefix = "CMP"
    var lastCode string

    err := tx.Table("company_models").
        Select("id_company").
        Order("id_company DESC").
        Limit(1).
        Clauses(clause.Locking{Strength: "UPDATE"}). 
        Row().
        Scan(&lastCode)

    if err != nil {
        return prefix + "0001", nil
    }

    numStr := lastCode[len(prefix):]
    currentNum, err := strconv.Atoi(numStr)
    if err != nil {
        return "", err
    }

    nextNum := currentNum + 1
    newCode := fmt.Sprintf("%s%04d", prefix, nextNum)

    return newCode, nil
}

func (s *CompanyService) GetCompanyForUser(ctx context.Context, user *user.UserModel) (*dto.CompanyInfo, error) {

	// ===== VALIDASI USER =====
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	// ===== FALLBACK (optional tapi recommended) =====
	if user.CompanyID == nil && len(user.Employees) > 0 {
		user.CompanyID = &user.Employees[0].CompanyID
	}

	if user.CompanyID == nil {
		return nil, fmt.Errorf("user has no company")
	}

	// ===== GET COMPANY =====
	company, err := s.CompanyRepo.GetCompanyWithPackage(ctx, *user.CompanyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	var startDate, endDate, packageName string

	if len(company.CompanySubscriptionHistory) > 0 {
		p := company.CompanySubscriptionHistory[0]

		startDate = helpers.DateFormat(p.StartDate)

		if p.EndDate != nil {
			endDate = helpers.DateFormat(*p.EndDate)
		}

		if p.Subscription != nil {
			packageName = p.Subscription.Name
		}
	}

	return &dto.CompanyInfo{
		ID:            company.ID.String(),
		CompanyName:   company.Name,
		EmployeeCount: company.EmployeeCount,
		IsActive:      true,
		IsTrial:       false,
		StartDate:     startDate,
		EndDate:       endDate,
		PackageName:   packageName,
	}, nil
}


func (s *CompanyService) CreateDefaultCompany(ctx context.Context, name string) (*company.CompanyModel, error) {

	var result *company.CompanyModel

	err := s.CompanyRepo.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 🔍 1. Check existing (SAFE)
		existing, err := s.CompanyRepo.FindByName(tx, name)
		if err != nil {
			return err
		}
		if existing != nil {
			result = existing
			return nil
		}

		// 🔍 2. Get default subscription (NO FATAL)
		pkg, err := s.SubscriptionRepo.FindByNameSubscription(tx, "Basic")
		if err != nil {
			return fmt.Errorf("failed to get default package: %w", err)
		}

		// 🔢 3. Generate company code
		code, err := s.generateCompanyCode(tx)
		if err != nil {
			return fmt.Errorf("failed generate company code: %w", err)
		}

		now := time.Now()
		trialEnd := now.AddDate(0, 0, pkg.DurationDays)

		// 🏗 4. Create company
		companyModel := &company.CompanyModel{
			ID:                    uuid.New(),
			IDCompany:             &code,
			Name:                  name,
			EmployeeCount:         0,
			CurrentSubscriptionID: &pkg.ID,
		}

		if err := tx.Create(companyModel).Error; err != nil {
			return err
		}

		// 📦 5. Create subscription history
		history := &subscription.SubscriptionHistoryModel{
			ID:             uuid.New(),
			CompanyID:      companyModel.ID,
			SubscriptionID: pkg.ID,
			StartDate:      now,
			EndDate:        &trialEnd,
			IsTrial:        true,
			IsActive:       true,
		}

		if err := s.SubscriptionRepo.CreateHistorySubscription(tx, history); err != nil {
			return err
		}

		result = companyModel
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}