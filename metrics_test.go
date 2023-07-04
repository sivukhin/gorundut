package gorundut

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeapAllocTotalBytes(t *testing.T) {
	total := GetHeapAllocTotal()
	require.Greater(t, total, uint64(0))
	t.Log(total)
}
