package api

import (
	"bytes"
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

	// testCases := []struct {
	// 	name          string
	// 	accoutID      int64
	// 	buildStubs    func(store *mockdb.MockStore)
	// 	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	// }{
	// 	{
	// 		name:      "OK",
	// 		accountID: account.ID,
	// 		buildStubs: func(store *mockdb.MockStore) {
	// 			store.EXPECT().
	// 				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
	// 				Times(1).
	// 				Return(account, nil)
	// 		},
	// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
	// 			require.Equal(t, http.StatusOK, recorder.Code)
	// 			requireBodyMatchAccount(t, recorder.Body, account)
	// 		},
	// 	},
	// 	// TODO
	// }

	ctrl := gomock.NewController(t) // Create a controler for later creating new mock store
	defer ctrl.Finish()             // The controller will check to see if all methods that were expected to be called were called

	store := mockdb.NewMockStore(ctrl) // Create new mock store

	// Build stubs for mock store, which we only care about GetAcount(), which is the only methods will be used by /accounts handler
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).            // How many times does this func should be called
		Return(account, nil) // Tell gomock to return some values (account value & nil err), this should match with GetAccount() in querier.go

	// Start to test server, and send HTTP req
	server := NewServer(store)
	recorder := httptest.NewRecorder() // Use httptest.Recorder (instead of starting a real HTTP server for testing an HTTP API), which can record the res of API req

	url := fmt.Sprintf("/accounts/%d", account.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil) // The req body is nil, cuz it's GET req
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req) // Send our API req thru the server router, and reocerder records it

	require.Equal(t, http.StatusOK, recorder.Code)
	// requireBodyMatchAccount(t, recorder.Body, account)
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
