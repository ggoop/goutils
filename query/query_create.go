package query

import (
	"fmt"

	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func SaveQuery(repo *repositories.MysqlRepo, item Query) (*Query, error) {
	if repo == nil || item.Code == "" {
		return nil, fmt.Errorf("参数不正确")
	}
	old := Query{}
	repo.Where("code=?", item.Code).Take(&old)
	if old.ID != "" {
		item.ID = old.ID
	} else {
		item.ID = utils.GUID()
		repo.Create(&item)
	}
	if item.Wheres != nil && len(item.Wheres) > 0 {
		count := 0
		for _, d := range item.Wheres {
			count = 0
			if repo.Model(d).Where("query_id=? and field=?", item.ID, d.Field).Count(&count); count == 0 {
				d.ID = utils.GUID()
				d.QueryID = item.ID
				repo.Create(&d)
			}
		}
	}
	return nil, nil
}
