package reg

import (
	"strings"

	"github.com/ggoop/goutils/md"
)

type RegObject struct {
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Addrs   []string `json:"addrs"`
	Content string   `json:"content"`
	Time    *md.Time `json:"time"`
}

func (s RegObject) Key() string {
	return strings.ToLower(s.Code)
}
