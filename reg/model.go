package reg

import (
	"strings"

	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"
)

type RegObject struct {
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Addrs   []string      `json:"addrs"`
	Content string        `json:"content"`
	Time    *md.Time      `json:"time"`
	Configs *utils.Config `json:"configs"`
}

func (s RegObject) Key() string {
	return strings.ToLower(s.Code)
}
