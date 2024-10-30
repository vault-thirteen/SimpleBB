package nt

import (
	"fmt"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// A list of commonly used notification template names.
const (
	Name_R    = Name("R")
	Name_FT   = Name("FT")
	Name_FRT  = Name("FRT")
	Name_FUT  = Name("FUT")
	Name_FUMT = Name("FUMT")
	Name_FMTU = Name("FMTU")
)

const (
	ErrF_ComponentCount = "number of template components is wrong: %v"
)

// Name is a notification template name.
type Name = cmb.Text

// checkComponentsCount ensures that each component is included no more than
// once.
func checkComponentsCount(name string) (err error) {
	var cc = make(map[rune]cmb.Count)
	var isFound bool
	for _, symbol := range ([]rune)(name) {
		_, isFound = cc[symbol]
		if !isFound {
			cc[symbol] = 0
		}
		cc[symbol]++
	}

	for symbol, count := range cc {
		if count > 1 {
			return fmt.Errorf(ErrF_ComponentCount, string(symbol))
		}
	}

	return nil
}
