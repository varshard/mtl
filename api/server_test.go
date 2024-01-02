//go:build integration

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/varshard/mtl/api/handlers"
	"github.com/varshard/mtl/api/handlers/responses"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"github.com/varshard/mtl/tests"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

const Authorization = "Authorization"

var (
	server      *httptest.Server
	db          *gorm.DB
	token       string
	bearerToken string
)

func TestMain(m *testing.M) {
	var err error
	conf := config.ReadEnv()
	db, err = database.InitDB(&conf.DBConfig)
	if err != nil {
		fmt.Println("fail to connect to the test database")
		return
	}
	defer tests.Truncate(db)

	r := MTLServer{}.InitRoutes(db, &conf)
	server = httptest.NewServer(r)
	defer server.Close()

	token, err = login("test_user", "password")
	if err != nil {
		fmt.Println("fail to acquire an authentication token")
		return
	}
	bearerToken = "Bearer " + token

	m.Run()
}

func TestLogin(t *testing.T) {
	url := server.URL + "/login"
	test := []struct {
		name           string
		payload        handlers.LoginRequest
		assertFunc     func(t *testing.T, resp handlers.LoginResponse)
		expectedStatus int
	}{
		{
			name: "given a correct credentials, it should returns a token",
			payload: handlers.LoginRequest{
				Username: "test",
				Password: "password",
			},
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T, resp handlers.LoginResponse) {
				assert.NotEmpty(t, resp.Token)
			},
		},
		{
			name: "given an incorrect credentials, it should returns a token",
			payload: handlers.LoginRequest{
				Username: "test",
				Password: "incorrect password",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			rawPayload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(rawPayload))
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			loginResp := handlers.LoginResponse{}
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&loginResp))

			if tt.assertFunc != nil {
				tt.assertFunc(t, loginResp)
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	url := server.URL + "/vote_items"
	defer tests.Truncate(db)

	payload := database.VoteItem{
		Name:        "test item",
		Description: "description",
	}

	test := []struct {
		name           string
		token          string
		payload        database.VoteItem
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "it should create an item successfully",
			token:          bearerToken,
			payload:        payload,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "it should returns 401 if the token missing Bearer prefix",
			token:          token,
			payload:        payload,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			rawPayload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawPayload))
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			loginResp := handlers.LoginResponse{}
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&loginResp))
		})
	}
}
func TestUpdate(t *testing.T) {
	url := server.URL + "/vote_items"
	tests.SeedDB(db)
	defer tests.Truncate(db)

	test := []struct {
		name           string
		token          string
		id             uint
		payload        database.VoteItem
		assertFunc     func(t *testing.T)
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "it should update an item successfully",
			token: bearerToken,
			id:    1,
			payload: database.VoteItem{
				Name:        "updated item",
				Description: "updated description",
			},
			assertFunc: func(t *testing.T) {
				updated := database.VoteItem{}
				assert.NoError(t, db.Table(database.TableVoteItem).Where("id = ?", 1).Take(&updated).Error)

				assert.Equal(t, "updated item", updated.Name)
				assert.Equal(t, "updated description", updated.Description)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "it should returns 404 ",
			id:             0,
			token:          bearerToken,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			rawPayload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", url, tt.id), bytes.NewBuffer(rawPayload))
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.assertFunc != nil {
				tt.assertFunc(t)
			}
		})
	}
}

func TestGetVoteItems(t *testing.T) {
	url := server.URL + "/vote_items"

	tests.SeedDB(db)
	defer tests.Truncate(db)

	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 3, VoteItemID: 3}).Error)
	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 2, VoteItemID: 3}).Error)
	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 1, VoteItemID: 1}).Error)

	test := []struct {
		name           string
		token          string
		assertFunc     func(t *testing.T, data []database.VoteItem)
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "it should returns items sorted by vote_count",
			token:          bearerToken,
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T, data []database.VoteItem) {
				assert.NotZero(t, data)
				prev := data[0]
				assert.NotZero(t, prev.VoteCount)
				for i := 1; i < len(data); i++ {
					item := data[i]
					assert.True(t, prev.VoteCount >= item.VoteCount)

					prev = data[i]
				}
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.assertFunc != nil {
				respBody := responses.DataResponse[[]database.VoteItem]{}
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
				tt.assertFunc(t, respBody.Data)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	url := server.URL + "/vote_items"
	tests.SeedDB(db)
	defer tests.Truncate(db)

	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 3, VoteItemID: 3}).Error)

	test := []struct {
		name           string
		token          string
		id             uint
		assertFunc     func(t *testing.T)
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "it should delete an item successfully",
			token: bearerToken,
			id:    1,
			assertFunc: func(t *testing.T) {
				var count int64
				assert.NoError(t, db.Table(database.TableVoteItem).Where("id = ?", 1).Count(&count).Error)
				assert.Zero(t, count)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "it should returns bad request if the item has been voted",
			token:          bearerToken,
			id:             3,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "it should returns 404 ",
			id:             0,
			token:          bearerToken,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", url, tt.id), nil)
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.assertFunc != nil {
				tt.assertFunc(t)
			}
		})
	}
}
func TestVote(t *testing.T) {
	url := server.URL + "/vote_items"
	tests.SeedDB(db)
	defer tests.Truncate(db)

	userVoteRepo := repository.VoteRepository{DB: db}
	token2, err := login("test", "password")
	assert.NoError(t, err)

	bearerToken2 := "Bearer " + token2

	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 3, VoteItemID: 3}).Error)

	test := []struct {
		name           string
		token          string
		id             uint
		assertFunc     func(t *testing.T)
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "it should vote for the specified item",
			token: bearerToken2,
			id:    1,
			assertFunc: func(t *testing.T) {
				voted, err := userVoteRepo.HasVote(1, 1)
				assert.NoError(t, err)
				assert.True(t, voted)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "it should returns bad request if user has already voted",
			token:          bearerToken,
			id:             3,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "it should returns 404 if the item doesn't exist",
			id:             0,
			token:          bearerToken2,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%d/vote", url, tt.id), bytes.NewBuffer(nil))
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.assertFunc != nil {
				tt.assertFunc(t)
			}
		})
	}
}

func TestResetItems(t *testing.T) {
	url := server.URL + "/vote_items/reset"
	tests.SeedDB(db)
	defer tests.Truncate(db)

	assert.NoError(t, db.Table(database.TableUserVote).Create(&database.UserVote{UserID: 3, VoteItemID: 3}).Error)

	test := []struct {
		name           string
		token          string
		assertFunc     func(t *testing.T)
		expectedStatus int
	}{
		{
			name:           "it should returns 401 if the bearer token is invalid",
			token:          "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "it should truncate vote items successfully",
			token:          bearerToken,
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T) {
				var count int64
				assert.NoError(t, db.Table(database.TableVoteItem).Count(&count).Error)
				assert.Zero(t, count)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(nil))
			assert.NoError(t, err)
			req.Header.Set(Authorization, tt.token)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.assertFunc != nil {
				tt.assertFunc(t)
			}
		})
	}
}

func login(username, password string) (string, error) {
	rawPayload, err := json.Marshal(handlers.LoginRequest{
		Username: username,
		Password: password,
	})

	resp, err := http.DefaultClient.Post(server.URL+"/login", "application/json", bytes.NewBuffer(rawPayload))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	loginResp := handlers.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return "", err
	}

	return loginResp.Token, nil
}
