package repository

import (
	"errors"
	"github.com/varshard/mtl/infrastructure/database"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r UserRepository) FindUser(name string) (*database.User, error) {
	u := &database.User{}
	err := r.DB.Table(database.TableUser).Select("id, name, password").
		Where("name = ?", name).Limit(1).Take(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xErr.NewErrNotFound(errors.New("user not found"))
	} else if err != nil {
		return nil, err
	}

	return u, nil
}
