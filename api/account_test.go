package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/kvnyijia/bank-app/db/mock"
	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Build stubs for mock store, which we only care about GetAcount(), which is the only methods will be used by /accounts handler
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).            // How many times does this func should be called
					Return(account, nil) // Tell gomock to return some values (account value & nil err), this should match with GetAccount() in querier.go
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Build stubs for mock store, which we only care about GetAcount(), which is the only methods will be used by /accounts handler
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                           // How many times does this func should be called
					Return(db.Account{}, sql.ErrNoRows) // This should match with GetAccountRequest() in api/account.go
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// Build stubs for mock store, which we only care about GetAcount(), which is the only methods will be used by /accounts handler
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                             // How many times does this func should be called
					Return(db.Account{}, sql.ErrConnDone) // This should match with GetAccountRequest() in api/account.go
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0, // Cuz the min ID should be 1
			buildStubs: func(store *mockdb.MockStore) {
				// Build stubs for mock store, which we only care about GetAcount(), which is the only methods will be used by /accounts handler
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) // How many times does this func should be called
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // Create a controler for later creating new mock store
			defer ctrl.Finish()             // The controller will check to see if all methods that were expected to be called were called

			store := mockdb.NewMockStore(ctrl) // Create new mock store

			// Build stubs
			tc.buildStubs(store)

			// Start to test server, and send HTTP req
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder() // Use httptest.Recorder (instead of starting a real HTTP server for testing an HTTP API), which can record the res of API req

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil) // The req body is nil, cuz it's GET req
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req) // Send our API req thru the server router, and reocerder records it
			tc.checkResponse(t, recorder)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
