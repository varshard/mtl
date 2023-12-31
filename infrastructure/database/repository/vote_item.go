package repository

import (
	"errors"
	"github.com/varshard/mtl/infrastructure/database"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"gorm.io/gorm"
)

const (
	ErrVoteItemNotFound = "vote item not found"
	ErrRemoveVotedItem  = "can't remove the vote item"
)

type ItemRepository struct {
	DB *gorm.DB
}

func (r ItemRepository) Exist(id uint) (bool, error) {
	var count int64 = 0
	err := r.DB.Table(database.TableVoteItem).Count(&count).Where("id = ?", id).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r ItemRepository) GetItems() ([]database.VoteItem, error) {
	items := make([]database.VoteItem, 0)
	err := r.DB.Table("vote_item v").Select("id, name, description, vote_count").
		Joins("LEFT JOIN (SELECT vote_item_id, COUNT(*) AS vote_count FROM user_vote GROUP_BY vote_item_id) u ON u.vote_item_id = v.id").
		Order("vote_count DESC").Find(&items).Error
	if err != nil {
		return items, err
	}

	return items, nil
}

func (r ItemRepository) Create(item database.VoteItem) (*database.VoteItem, error) {
	err := r.DB.Table(database.TableVoteItem).Create(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r ItemRepository) Update(id uint, item database.VoteItem) error {
	exist := &database.VoteItem{ID: id}

	err := r.DB.Table(exist.TableName()).Where(exist).Take(&exist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return xErr.NewErrNotFound(errors.New(ErrVoteItemNotFound))
	} else if err != nil {
		return err
	}

	if err := r.DB.Table(item.TableName()).Where(exist).Updates(item).Error; err != nil {
		return err
	}

	return nil
}

func (r ItemRepository) Remove(id uint) error {
	item := &database.VoteItem{}
	err := r.DB.Table("vote_item v").Select("id, name, description, vote_count").
		Joins("LEFT JOIN (SELECT vote_item_id, COUNT(*) AS vote_count FROM user_vote GROUP_BY vote_item_id) u ON u.vote_item_id = v.id").
		Where("v.id = ?", id).
		First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return xErr.NewErrNotFound(errors.New(ErrVoteItemNotFound))
	} else if err != nil {
		return err
	}
	if item.VoteCount > 0 {
		return xErr.NewErrInvalidInput(errors.New(ErrRemoveVotedItem))
	}

	return nil
}

func (r ItemRepository) ResetItems() error {
	if err := r.DB.Raw("TRUNCATE user_vote").Error; err != nil {
		return err
	}
	return r.DB.Raw("TRUNCATE vote_item").Error
}
