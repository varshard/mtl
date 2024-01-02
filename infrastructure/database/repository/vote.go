package repository

import (
	"github.com/varshard/mtl/infrastructure/database"
	"gorm.io/gorm"
)

type VoteRepository struct {
	DB *gorm.DB
}

func (r VoteRepository) ResetVotes() error {
	return r.DB.Exec("DELETE FROM user_vote").Error
}

func (r VoteRepository) HasVote(itemID, userID uint) (bool, error) {
	var count int64 = 0
	if err := r.DB.Table(database.TableUserVote).Where("user_id = ? AND vote_item_id = ?", userID, itemID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r VoteRepository) IsVotable(userID uint) (bool, error) {
	var count int64 = 0
	if err := r.DB.Table(database.TableUserVote).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r VoteRepository) Vote(itemID, userID uint) error {
	v := database.UserVote{
		VoteItemID: itemID,
		UserID:     userID,
	}
	return r.DB.Table(v.TableName()).Create(&v).Error
}

func (r VoteRepository) ClearVote(itemID uint) error {
	v := database.UserVote{
		VoteItemID: itemID,
	}
	return r.DB.Table(v.TableName()).Where(v).Delete(&v).Error
}
