package nat

import (
	"github.com/blockchainservice/common"
)

var log common.Logger

func init() {
	DisableLog()
}

func DisableLog() {
	log = common.Disabled
}

func UseLogger(logger common.Logger) {
	log = logger
}
