package main

import (
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/query"
	"github.com/ggoop/goutils/repositories"
)

func main() {
	mysql := repositories.NewMysqlRepo()
	md.Migrate(mysql, &query.Query{}, &query.QueryField{})
}
