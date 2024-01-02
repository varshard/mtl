package database

const (
	TableUser     = "user"
	TableVoteItem = "vote_item"
	TableUserVote = "user_vote"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string
}

func (User) TableName() string {
	return "user"
}

type VoteItem struct {
	ID          uint   `json:"-" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VoteCount   uint   `json:"vote_count" gorm:"<-:false"`
}

func (VoteItem) TableName() string {
	return TableVoteItem
}

type UserVote struct {
	VoteItemID uint `json:"vote_item_id"`
	UserID     uint `json:"user_id"`
}

func (UserVote) TableName() string {
	return TableUserVote
}
