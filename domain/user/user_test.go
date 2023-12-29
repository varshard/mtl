//go:build unit

package user_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/varshard/mtl/domain/user"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"testing"
)

func TestUserRepository(t *testing.T) {
	conf := config.ReadEnv()
	db, err := database.InitDB(&conf.DBConfig)

	assert.NoError(t, err)
	repo := user.Repository{DB: db}

	t.Run("FindUser", func(t *testing.T) {
		t.Run("should return a user matching the user name", func(t *testing.T) {
			u, err := repo.FindUser("test")

			require.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, "test", u.Name)
			assert.NotZero(t, u.ID)
		})

		t.Run("should return nil if user doesn't exist", func(t *testing.T) {
			u, err := repo.FindUser("not exist")

			assert.EqualError(t, err, "user not found")
			assert.True(t, errors.As(err, &xErr.ErrNotFound{}))
			assert.Nil(t, u)
		})
	})
}
