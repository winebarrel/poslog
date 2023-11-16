package utils_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/poslog/utils"
)

func TestReadLine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString("select 1\n") //nolint:errcheck
	f.WriteString("select 2\n") //nolint:errcheck
	f.WriteString("select 3\n") //nolint:errcheck
	f.Sync()                    //nolint:errcheck
	f.Seek(0, io.SeekStart)     //nolint:errcheck

	buf := bufio.NewReader(f)

	for _, expected := range []string{"select 1\n", "select 2\n", "select 3\n"} {
		line, err := utils.ReadLine(buf)
		require.NoError(err)
		assert.Equal(expected, line)
	}
}
