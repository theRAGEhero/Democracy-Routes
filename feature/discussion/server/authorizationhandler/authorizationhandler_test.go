package authorizationhandler_test

import (
	"testing"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/require"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
)

func TestAuthorizationHandler(t *testing.T) {
	t.Parallel()

	ctx := testhelper.Context(t)

	ofgaClient := testhelper.TmpOfgaClient(t)

	body := client.ClientWriteAuthorizationModelRequest{
		SchemaVersion: "1.1",
		TypeDefinitions: []openfga.TypeDefinition{
			{Type: "user", Relations: &map[string]openfga.Userset{}},
			{
				Type: "document",
				Relations: &map[string]openfga.Userset{
					"writer": {
						This: &map[string]interface{}{},
					},
					"viewer": {Union: &openfga.Usersets{
						Child: []openfga.Userset{
							{This: &map[string]interface{}{}},
							{ComputedUserset: &openfga.ObjectRelation{
								Object:   openfga.PtrString(""),
								Relation: openfga.PtrString("writer"),
							}},
						},
					}},
				},
				Metadata: &openfga.Metadata{
					Relations: &map[string]openfga.RelationMetadata{
						"writer": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"viewer": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
					},
				},
			}},
	}

	data, err := ofgaClient.WriteAuthorizationModel(ctx).Body(body).Execute()
	require.NoError(t, err, "creating model")

	t.Log(data.AuthorizationModelId)

	// TODO: save this model to OpenFGA DSL as a file and recreate from the saved file.
}
