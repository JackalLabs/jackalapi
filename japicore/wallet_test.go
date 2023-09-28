package japicore

import (
	"testing"
)

func TestInitWalletSession(t *testing.T) {
	t.Setenv("JAPI_SEED", "slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum")

	InitWalletSession()
}
