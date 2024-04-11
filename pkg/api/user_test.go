package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/aseerkt/go-simple-bank/pkg/mockdb"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func createRandomUser(t *testing.T) (db.User, string) {
	password := gofakeit.Password(true, true, true, true, true, 10)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	require.NoError(t, err)

	return db.User{
		Username:       gofakeit.Username(),
		FullName:       gofakeit.Name(),
		Email:          gofakeit.Email(),
		HashedPassword: string(hashedPassword),
	}, password
}

func TestCreateUserAPI(t *testing.T) {

	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusCreated)

				var result db.User
				err := json.Unmarshal(recorder.Body.Bytes(), &result)

				require.NoError(t, err)
				require.NotEmpty(t, result)

				require.Equal(t, result.Username, user.Username)
				require.Equal(t, result.FullName, user.FullName)
				require.Equal(t, result.Email, user.Email)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "@$@!@$!@#",
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(store)

			url := "/users"

			server.LoadRoutes()

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestLoginUser(t *testing.T) {

	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(user, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusOK)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": ")(!@(#))",
				"password": password,
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetUser(gomock.Any(), ")(!@(#))").Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusBadRequest)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusNotFound)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusInternalServerError)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username": user.Username,
				"password": "custompassword",
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(user, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusUnauthorized)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(store)
			server.LoadRoutes()

			tc.buildStubs(store)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}
