package database

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain"
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/jinzhu/gorm"
)

// シーケンスのRepository
type SequenceRepository struct{}

// 指定したイベントのシーケンスをインクリメントし、値を返却します。
func (rep SequenceRepository) Increment(db *gorm.DB, event enum.Event) int {
	tx := db.Begin()
	var currentSequence domain.Sequence
	tx.Where("event = ?", event).Find(&currentSequence)
	next := currentSequence.Sequence + 1
	tx.Model(&currentSequence).Update("sequence", next)
	tx.Commit()
	return next
}
