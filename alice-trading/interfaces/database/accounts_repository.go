package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/jinzhu/gorm"
)

// アカウントのRepository
type AccountsRepository struct{}

func (rep AccountsRepository) FindByID(db *gorm.DB, id int) domain.Accounts {
	var account domain.Accounts
	db.Find(&account, id)
	return account
}

func (rep AccountsRepository) Update(db *gorm.DB, params map[string]interface{}) domain.Accounts {
	tx := db.Begin()
	var account domain.Accounts
	tx.First(&account)
	tx.Model(&account).Updates(params)
	tx.First(&account)
	tx.Commit()
	return account
}
