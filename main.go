package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 10; i++ {
		glog.Info(utils.GUID())
	}
	testOracle()
}
func testOracle() {
	repo := repositories.Default()
	//md.Migrate(repo)
	//md.InitMD_Completed = true
	m := &testTable{}
	//md.Migrate(repo, m)
	count := 0
	repo.Create(m)

	mList := make([]testTable, 0)
	repo.Find(mList, "is_system=?", 0)

	//take
	m = &testTable{}
	repo.Model(testTable{}).Take(m, "is_system=?", 0)

	glog.Error(repo.Model(testTable{}).Where("Is_System=?", 1).Count(&count).Error)
	//
	glog.Info(utils.GUID())
}

type testTable struct {
	ID       string `gorm:"primary_key;size:50" json:"id"`
	IsSystem int
}

func (s testTable) MD() *md.Mder {
	return &md.Mder{ID: "test.table", Domain: "test", Name: "标签"}
}

func (s testTable) TableName() string {
	return "testTable"
}
