package gorundut

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	directory, err := os.MkdirTemp("", "test-dump-")
	require.Nil(t, err)
	defer os.RemoveAll(directory)

	prefix := time.Now().Format("2006-01-02_15-04-05")
	dump, err := MakeProfilesDump(path.Join(directory, prefix), AllProfileNames...)
	require.Nil(t, err)
	t.Log(dump)
}
