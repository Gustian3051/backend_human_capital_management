package database

import (
	"backend/internal/domain/attendance"
	"backend/internal/domain/common"
	"backend/internal/domain/company"
	"backend/internal/domain/department"
	"backend/internal/domain/employee"
	"backend/internal/domain/package"
	"backend/internal/domain/position"
	"backend/internal/domain/role_permission"
	"backend/internal/domain/shift"
	"backend/internal/domain/user"
	"backend/internal/domain/work_leave"
	"backend/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// daftar semua model
func getModels() []interface{} {
	return []interface{}{
		// package
		&pkg.PackageModel{},
		&pkg.PackageHistoryModel{},
		
		// company
		&company.CompanyModel{},

		// department
		&department.DepartmentModel{},

		// position
		&position.PositionModel{},

		// role permission
		&role_permission.RoleModel{},
		&role_permission.PermissionModel{},
		&role_permission.RolePermissionModel{},

		// user
		&user.UserModel{},

		// employee
		&employee.EmployeeModel{},

		// shift
		&shift.ShiftModel{},
		&shift.WorkDayModel{},
		&shift.ShiftWorkDayModel{},
		&shift.ShiftAssignmentModel{},

		// attendance
		&attendance.AttendanceRulesModel{},
		&attendance.AttendanceModel{},

		// work leave
		&work_leave.LeaveBalanceModel{},
		&work_leave.WorkLeaveModel{},

		// common
		&common.LogActionModel{},
		&common.CasbinRule{},
	}
}

// =============================
// SAFE MIGRATE (DEFAULT)
// =============================
func AutoMigrate(db *gorm.DB) error {
	logger.Log.Info("Running auto migration (safe mode)...")

	return db.AutoMigrate(getModels()...)
}

// =============================
// FORCE RESET MIGRATE 🔥
// =============================
func ResetMigrate(db *gorm.DB) error {
	logger.Log.Warn("⚠️ RESET MIGRATION: dropping all tables...")

	err := db.Migrator().DropTable(getModels()...)
	if err != nil {
		logger.Log.Error("Failed to drop tables", zap.Error(err))
		return err
	}

	logger.Log.Info("Recreating tables...")

	err = db.AutoMigrate(getModels()...)
	if err != nil {
		logger.Log.Error("Failed to migrate tables", zap.Error(err))
		return err
	}

	logger.Log.Info("Migration completed successfully")

	return nil
}