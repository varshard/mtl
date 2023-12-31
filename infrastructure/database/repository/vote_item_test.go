//go:build integration

package repository

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/varshard/mtl/domain/vote"
	"github.com/varshard/mtl/infrastructure/database"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"testing"
)

const TestUserID = 3

func TestVoteItemRepository(t *testing.T) {
	repo := ItemRepository{DB: db}

	t.Run("Exist", func(t *testing.T) {
		tests := []struct {
			name     string
			id       uint
			expected bool
			err      string
		}{
			{
				name:     "should returns true if the vote item exist",
				id:       1,
				expected: true,
			},
			{
				name:     "should returns false if the vote item doesn't exist",
				id:       99,
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual, err := repo.Exist(tt.id)

				if tt.err != "" {
					assert.EqualError(t, err, tt.err)
				} else {
					assert.Equal(t, tt.expected, actual)
				}
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		defer tearDownVoteItem()

		tests := []struct {
			name  string
			input database.VoteItem
			err   any
		}{
			{
				name: "should create an item successfully",
				input: database.VoteItem{
					CreatedBy:   TestUserID,
					Name:        "test item",
					Description: "description",
				},
			},
			{
				name: "should returns an error if the creator ID is invalid",
				input: database.VoteItem{
					CreatedBy:   99,
					Name:        "test item 2",
					Description: "should returns an error",
				},
				err: xErr.NewErrUnexpected(errors.New("dummy")),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				item, err := repo.Create(tt.input)

				if tt.err != nil {
					assert.True(t, errors.As(err, &tt.err))
				} else {
					assert.NotZero(t, item.ID)

					input := tt.input
					assert.Equal(t, input.Name, item.Name)
					assert.Equal(t, input.Description, item.Description)
					assert.Equal(t, input.CreatedBy, item.CreatedBy)
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDownVoteItem()

		item := &database.VoteItem{
			CreatedBy:   TestUserID,
			Name:        "test create",
			Description: "description",
		}
		require.NoError(t, db.Create(&item).Error)

		tests := []struct {
			name  string
			id    uint
			input vote.UpdateVoteItem
			err   any
		}{
			{
				name: "should update an item successfully",
				id:   item.ID,
				input: vote.UpdateVoteItem{
					Name:        "new name",
					Description: "new description",
				},
			},
			{
				name: "should returns an error if the ID is invalid",
				id:   99,
				input: vote.UpdateVoteItem{
					Name:        "this should fail",
					Description: "should returns an error",
				},
				err: xErr.NewErrNotFound(errors.New("dummy")),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.Update(tt.id, tt.input)

				if tt.err != nil {
					assert.True(t, errors.As(err, &tt.err))
				} else {
					actual := &database.VoteItem{}
					assert.NoError(t, db.Table(database.TableVoteItem).Where("id = ?", tt.id).Take(&actual).Error)
					assert.Equal(t, tt.id, actual.ID)

					input := tt.input
					assert.Equal(t, input.Name, actual.Name)
					assert.Equal(t, input.Description, actual.Description)
					assert.Equal(t, item.CreatedBy, actual.CreatedBy)
				}
			})
		}
	})
}

func tearDownVoteItem() {
	db.Exec("DELETE FROM vote_item WHERE created_by = ?", TestUserID)
}
