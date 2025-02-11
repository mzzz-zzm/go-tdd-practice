package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type MeanyGreeter interface {
	Curse(name string) (string, error)
}

func CurseSpecifications(t testing.TB, meany MeanyGreeter) {
	got, err := meany.Curse("Chris")
	assert.NoError(t, err)
	assert.Equal(t, got, "no way, Chris")
}
