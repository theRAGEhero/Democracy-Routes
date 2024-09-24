package cli_test

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

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
				Args: []string{"create", "user", "-name=Dima", "-pass="},
				Out:  &buf,
				DB:   testhelper.TmpDB(t),
			}),
			"creating user",
		)

		// Then the user Dima exists.

		assert.Contains(t, buf.String(), `"ID":"`)
		assert.Contains(t, buf.String(), `"Name":"Dima"`)

		ps := strings.Index(buf.String(), `"Password":"`) + len(`"Password":"`)
		pe := strings.Index(buf.String()[ps:], `"`) + ps
		assert.Equal(t, 16, utf8.RuneCountInString(buf.String()[ps:pe]))
	})
}
