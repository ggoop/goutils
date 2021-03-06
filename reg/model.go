package reg

import (
	"strings"

	"github.com/ggoop/goutils/utils"
)

type RegObject struct {
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	PublicAddress string        `json:"public_address"`
	Content       string        `json:"content"`
	Time          *utils.Time   `json:"time"`
	Configs       *utils.Config `json:"configs"`
}

func (s RegObject) Key() string {
	return strings.ToLower(s.Code)
}
