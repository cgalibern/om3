package usergroup

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUidFromS(t *testing.T) {
	cases := map[string]uint32{
		"root": uint32(0),
		"0":    uint32(0),
	}
	for s, expected := range cases {
		t.Run("with valid user: "+s, func(t *testing.T) {
			id, err := UIDFromString(s)
			require.Nil(t, err)
			assert.Equal(t, expected, id)
		})
	}
	t.Run("with invalid user", func(t *testing.T) {
		_, err := UIDFromString("badUserX")
		require.NotNil(t, err)
	})
}

func TestGidFromS(t *testing.T) {
	cases := map[string]uint32{
		"daemon": uint32(1),
		"1":      uint32(1),
	}
	if runtime.GOOS == "solaris" {
		cases["daemon"] = uint32(12)
	}
	for s, expected := range cases {
		t.Run("valid group: "+s, func(t *testing.T) {
			id, err := GIDFromString(s)
			require.Nil(t, err)
			assert.Equal(t, expected, id)
		})
	}
	t.Run("invalid group", func(t *testing.T) {
		_, err := GIDFromString("badGroupY")
		require.NotNil(t, err)
	})
}
