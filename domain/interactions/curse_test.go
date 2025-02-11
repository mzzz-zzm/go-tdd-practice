package interactions

import (
	"testing"

	"github.com/mzzz-zzm/go-tdd-practice/specifications"
)

func TestCurse(t *testing.T) {
	specifications.CurseSpecifications(t, specifications.CurseAdapter(Curse))
}
