package configs

import (
	"github.com/ggoop/goutils/utils"
)

var Default *utils.Config

func init() {
	Default = utils.DefaultConfig
}
