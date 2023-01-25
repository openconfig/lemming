package bgp

import (
	"fmt"

	"github.com/wenovus/gobgp/v3/pkg/log"
)

const (
	debugBGP = false
)

func debugBGPPrint(l log.Logger, s string) {
	if debugBGP {
		l.Info(fmt.Sprintf("DEBUG/BGP %s", s), log.Fields{"Topic": "config"})
	}
}
