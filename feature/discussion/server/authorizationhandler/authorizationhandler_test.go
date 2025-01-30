package authorizationhandler_test

import (
	_ "embed"
	"testing"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
	"github.com/stretchr/testify/assert"
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

	// TODO: read model from OpenFGA DSL string.
	// It looks like OpenFGA language parsed can't convert from native model to OpenFGA DSL string rught now.
}

//go:embed model.fga
var dslModel string

func TestLanguage(t *testing.T) {
	t.Parallel()

	generatedJsonString, err := transformer.TransformDSLToJSON(dslModel)
	require.NoError(t, err, "dsl > json")

	generatedDsl, err := transformer.TransformJSONStringToDSL(generatedJsonString)
	require.NoError(t, err, "json > dsl")

	// generatedProto, err := transformer.TransformDSLToProto(dslString)
	// require.NoError(t, err, "transforming from dsl to proto")

	// _, err = transformer.TransformJSONProtoToDSL(generatedProto)
	// require.NoError(t, err, "transforming from proto to dsl")

	assert.Equal(t, dslModel, *generatedDsl, "outcome is different")
}
