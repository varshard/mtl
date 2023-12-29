package vote

import (
	"gorm.io/gorm"
)

const TableVote = "user_vote"

type Vote struct {
	VoteItemID uint
	UserID     uint
}

type Repository struct {
	DB *gorm.DB
}

func (r Repository) ResetVote() error {
	return r.DB.Raw("TRUNCATE user_vote").Error
}

func (r Repository) Vote(itemID, userID uint) error {
	v := Vote{
		VoteItemID: itemID,
		UserID:     userID,
	}
	if err := r.DB.Table(TableVote).Create(&v).Error; err != nil {
		return err
	}
	return nil
}

func (r Repository) UnVote(itemID, userID uint) error {
	v := Vote{
		VoteItemID: itemID,
		UserID:     userID,
	}
	if err := r.DB.Table(TableVote).Delete(&v).Error; err != nil {
		return err
	}
	return nil
}
