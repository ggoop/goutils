package reg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"
	"github.com/golang/glog"
	"io"
	"os"
	"path"
)

var store *RegStoreSv

type RegStoreSv struct {
	data   map[string]*RegObject
	dbFile string
}

func NewRegStoreSv() *RegStoreSv {
	return &RegStoreSv{data: make(map[string]*RegObject)}
}

func (s *RegStoreSv) Register() {
	addrs := make([]string, 0)
	if host := configs.Default.App.Address; host != "" {
		addrs = append(addrs, host)
	} else {
		ips := utils.GetIpAddrs()
		for _, item := range ips {
			addrs = append(addrs, fmt.Sprintf("http://%s:%s", item, configs.Default.App.Port))
		}
	}
	s.Add(RegObject{
		Code:  configs.Default.App.Code,
		Name:  configs.Default.App.Name,
		Addrs: addrs,
	})
}
func (s *RegStoreSv) Add(item RegObject) *RegObject {
	item.Time = md.NewTimePtr()
	s.data[item.Key()] = &item
	s.Store()
	return &item
}

func (s *RegStoreSv) Get(item RegObject) *RegObject {
	if old, ok := s.data[item.Key()]; ok {
		return old
	} else {
		return nil
	}
}

func (s *RegStoreSv) GetAll() []RegObject {
	items := make([]RegObject, 0)
	for _, item := range s.data {
		items = append(items, *item)
	}
	return items
}
func (s *RegStoreSv) Init() {
	if s.dbFile == "" {
		s.dbFile = utils.JoinCurrentPath(path.Join(configs.Default.App.Storage, "uploads", "regs"))
	}
	if !utils.PathExists(s.dbFile) {
		return
	}
	fi, err := os.Open(s.dbFile)
	if err != nil {
		glog.Error(err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		item := RegObject{}
		err = json.Unmarshal(a, &item)
		if err != nil {
			glog.Error(err)
			return
		}
		s.Add(item)
	}
}
func (s *RegStoreSv) Store() {
	items := s.GetAll()
	f, err := os.Create(s.dbFile)
	if err != nil {
		glog.Error(err)
		f.Close()
		return
	}
	for _, item := range items {
		b, err := json.Marshal(item)
		if err != nil {
			glog.Error(err)
			return
		}
		fmt.Fprintln(f, string(b))
		if err != nil {
			glog.Error(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		glog.Error(err)
		return
	}
}
