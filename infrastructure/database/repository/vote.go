package repository

import (
	"github.com/varshard/mtl/infrastructure/database"
	"gorm.io/gorm"
)

type VoteRepository struct {
	DB *gorm.DB
}

func (r VoteRepository) ResetVote() error {
	return r.DB.Raw("TRUNCATE user_vote").Error
}

func (r VoteRepository) HasVote(itemID, userID uint) (bool, error) {
	var count int64 = 0
	if err := r.DB.Table(database.TableUserVote).Where("user_id = ? AND vote_item_id = ?", userID, itemID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r VoteRepository) IsVoteable(itemID, userID uint) (bool, error) {
	var count int64 = 0
	if err := r.DB.Table(database.TableUserVote).Where("user_id = ? AND vote_item_id <> ?", userID, itemID).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r VoteRepository) Vote(itemID, userID uint) error {
	v := database.UserVote{
		VoteItemID: itemID,
		UserID:     userID,
	}
	if err := r.DB.Table(database.TableUserVote).Create(&v).Error; err != nil {
		return err
	}
	return nil
}

func (r VoteRepository) UnVote(itemID, userID uint) error {
	v := database.UserVote{
		VoteItemID: itemID,
		UserID:     userID,
	}
	if err := r.DB.Table(database.TableUserVote).Delete(&v).Error; err != nil {
		return err
	}
	return nil
}
