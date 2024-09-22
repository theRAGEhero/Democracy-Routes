package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
)

func TestCommandLineInterface(t *testing.T) {
	t.Parallel()

	t.Run("creating new user", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		// Given a new user Dima is added.

		require.NoError(
			t,
			cli.Run(common.Params{
				Args: []string{"create", "user", "-name=Dima", "-password=secret"},
				Out:  &buf,
				DB:   testhelper.TmpDB(t),
			}),
			"creating user",
		)

		// Then the user Dima exists.

		assert.Contains(t, buf.String(), "created")
	})
}
