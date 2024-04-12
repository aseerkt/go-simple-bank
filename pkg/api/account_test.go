package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/aseerkt/go-simple-bank/pkg/constants"
	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/aseerkt/go-simple-bank/pkg/mockdb"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createRandomAccount() (db.User, db.Account) {
	user := db.User{
		Username: gofakeit.Username(),
		FullName: gofakeit.Name(),
		Email:    gofakeit.Email(),
	}

	return user, db.Account{
		ID:       gofakeit.Int64(),
		Owner:    user.Username,
		Balance:  gofakeit.Int64(),
		Currency: gofakeit.RandomMapKey(constants.Currency).(string),
	}
}

func getAuthMiddleware(username string) func(t *testing.T, s *Server, r *http.Request) {
	return func(t *testing.T, s *Server, r *http.Request) {
		token, err := s.tokenMaker.CreateToken(username, 15*time.Minute)
		require.NoError(t, err)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func TestCreateAccountAPI(t *testing.T) {

	user, account := createRandomAccount()

	setupAuth := getAuthMiddleware(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, s *Server, r *http.Request)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			setupAuth: setupAuth,
			body: gin.H{
				"currency": "INR",
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    user.Username,
					Currency: "INR",
					Balance:  0,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, r.Code)
				var createdAccount db.Account
				err := json.Unmarshal(r.Body.Bytes(), &createdAccount)
				require.NoError(t, err)
				require.Equal(t, createdAccount, account)
			},
		},
		{
			name:      "InvalidCurrency",
			setupAuth: setupAuth,
			body: gin.H{
				"currency": "NONE",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name:      "UniqueViolation",
			setupAuth: setupAuth,
			body: gin.H{
				"currency": "INR",
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    user.Username,
					Currency: "INR",
					Balance:  0,
				}
				err := &pq.Error{
					Code: "23505",
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{}, err)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, r.Code)

			},
		},
		{
			name:      "InternalError",
			setupAuth: setupAuth,
			body: gin.H{
				"currency": "INR",
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    user.Username,
					Currency: "INR",
					Balance:  0,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := newTestServer(store)
			server.LoadRoutes()

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			require.NoError(t, err)

			tc.setupAuth(t, server, request)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	user, account := createRandomAccount()

	setupAuth := getAuthMiddleware(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, s *Server, r *http.Request)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			setupAuth: setupAuth,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)

				require.NoError(t, err)

				var gotAccount db.Account

				err = json.Unmarshal(data, &gotAccount)

				require.NoError(t, err)
				require.Equal(t, account, gotAccount)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: setupAuth,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: setupAuth,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: setupAuth,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(0)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "AccountMismatch",
			accountID: account.ID,
			setupAuth: getAuthMiddleware("alfred"),
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStub(store)

			server := newTestServer(store)
			server.LoadRoutes()

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, server, request)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {

	var username string
	var accounts []db.Account

	for range 5 {
		user, account := createRandomAccount()
		if username == "" {
			username = user.Username
		} else {
			account.Owner = username
		}
		accounts = append(accounts, account)
	}

	setPaginationQuery := func(r *http.Request, pageID int, pageSize int) {
		var query url.Values = make(url.Values)

		query.Add("page_id", fmt.Sprint(pageID))
		query.Add("page_size", fmt.Sprint(pageSize))

		r.URL.RawQuery = query.Encode()
	}

	testCases := []struct {
		name          string
		setQuery      func(r *http.Request)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			setQuery: func(r *http.Request) {
				setPaginationQuery(r, 1, 5)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  username,
					Offset: 0,
					Limit:  5,
				}
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
				var listedAccounts []db.Account
				err := json.Unmarshal(r.Body.Bytes(), &listedAccounts)
				require.NoError(t, err)
				require.Equal(t, listedAccounts, accounts)

			},
		},
		{
			name: "InvalidQuery",
			setQuery: func(r *http.Request) {
				setPaginationQuery(r, 0, 25)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name: "InternalServerError",
			setQuery: func(r *http.Request) {
				setPaginationQuery(r, 1, 5)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  username,
					Offset: 0,
					Limit:  5,
				}
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
	}

	setupAuth := getAuthMiddleware(username)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := newTestServer(store)
			server.LoadRoutes()

			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)

			require.NoError(t, err)

			setupAuth(t, server, request)

			tc.setQuery(request)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
