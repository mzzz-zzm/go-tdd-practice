package interactions

import (
	"testing"

	"github.com/mzzz-zzm/go-tdd-practice/specifications"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecifications(t, specifications.GreetAdapter(Greet))
}
