package usecase

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// シーケンスのRepository
type SequenceRepository interface {
	Increment(db *gorm.DB, event enum.Event) int
}
