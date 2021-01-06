package mastodon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAcct(t *testing.T) {
	t.Run("same-server", func(t *testing.T) {
		require.Equal(t, "name@domain.com", NormalizeAcct("https://domain.com", "name"))
		require.Equal(t, "name@domain.com", NormalizeAcct("https://domain.com", "name@domain.com"))
	})

	t.Run("diff-server", func(t *testing.T) {
		require.Equal(t, "name@other.com", NormalizeAcct("https://domain.com", "name@other.com"))
	})
}
