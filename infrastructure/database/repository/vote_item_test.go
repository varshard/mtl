//go:build integration

package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
}
