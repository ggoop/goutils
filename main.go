package main

import (
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		glog.Errorf("%v\r\n", utils.GUID())
	}
	// test_query()
}
func test_query() {
	mysql := repositories.NewMysqlRepo()
	tx := mysql.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.AutoMigrate(query.Query{})

	item := query.Query{}
	item.ID = utils.GUID()
	item.Name = "ddd"
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
	}

	item = query.Query{}
	item.ID = "01e94b9a020ac7178f2bfa163e9054be"
	item.Name = "ccc"
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
	}

	items := make([]query.Query, 0)
	if err := mysql.Find(&items).Error; err != nil {
		glog.Errorf("Find error is %s", err)
	}
	if err := tx.Commit().Error; err != nil {
		glog.Errorf("Commit error is %s", err)
	}
}
