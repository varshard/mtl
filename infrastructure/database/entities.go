package database

const (
	TableUser     = "user"
	TableVoteItem = "vote_item"
	TableUserVote = "user_vote"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Password string
}

func (User) TableName() string {
	return "user"
}

type VoteItem struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedBy   uint
	VoteCount   uint `gorm:"<-:false"`
}

func (VoteItem) TableName() string {
	return TableVoteItem
}

type UserVote struct {
	VoteItemID uint
	UserID     uint
}

func (UserVote) TableName() string {
	return TableUserVote
}
