package user

import (
	"errors"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Password string
}

func (User) TableName() string {
	return "user"
}

type Repository struct {
	DB *gorm.DB
}

func (r Repository) FindUser(name string) (*User, error) {
	u := &User{}
	err := r.DB.Table("user").Select("id, name, password").
		Where("name = ?", name).Limit(1).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xErr.NewErrNotFound(errors.New("user not found"))
	} else if err != nil {
		return nil, err
	}

	return u, nil
}
