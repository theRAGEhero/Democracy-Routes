package authorizationhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
)

func TestAuthorizationHandler(t *testing.T) {
	t.Parallel()

	ctx := testhelper.Context(t)

	ofgaClient := testhelper.TmpOfgaClient(t)

	storesRes, err := ofgaClient.ListStores(ctx).Execute()
	require.NoError(t, err, "listing stores")

	storeID, err := ofgaClient.GetStoreId()
	require.NoError(t, err, "getting client store")

	if assert.Len(t, storesRes.Stores, 1, "wrong number of stores") {
		assert.Equal(t, storeID, storesRes.Stores[0].Id, "got wrong store")
	}

	// TODO: create an authorization model.
}
