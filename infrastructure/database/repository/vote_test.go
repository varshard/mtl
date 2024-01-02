//go:build unit

package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/varshard/mtl/infrastructure/database"
	"testing"
)

func TestVoteRepository(t *testing.T) {
	repo := VoteRepository{DB: db}

	t.Run("Vote", func(t *testing.T) {
		db.Table(database.TableUserVote).Create(&database.UserVote{
			UserID:     1,
			VoteItemID: 2,
		})
		defer func() {
			assert.NoError(t, repo.ResetVotes())
		}()
		tests := []struct {
			name          string
			userID        uint
			itemID        uint
			expectedError bool
		}{
			{
				name:   "should create a record successfully",
				userID: TestUserID,
				itemID: 1,
			},
			{
				name:          "should returns an error if the user tried to vote again",
				userID:        1,
				itemID:        2,
				expectedError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.Vote(tt.itemID, tt.userID)
				if tt.expectedError {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("IsVotable", func(t *testing.T) {
		db.Table(database.TableUserVote).Create(&database.UserVote{
			UserID:     1,
			VoteItemID: 2,
		})

		defer repo.ResetVotes()

		tests := []struct {
			name          string
			userID        uint
			expected      bool
			expectedError bool
		}{
			{
				name:     "should returns true if the user hasn't vote on any item yet",
				userID:   TestUserID,
				expected: true,
			},
			{
				name:     "should returns false if the user has already casted their vote",
				userID:   1,
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual, err := repo.IsVotable(tt.userID)
				if tt.expectedError {
					assert.Error(t, err)
				}

				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}
