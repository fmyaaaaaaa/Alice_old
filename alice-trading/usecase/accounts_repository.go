package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// アカウントのRepository
type AccountRepository interface {
	FindByID(db *gorm.DB, id int) domain.Accounts
	Update(db *gorm.DB, params map[string]interface{}) domain.Accounts
}
