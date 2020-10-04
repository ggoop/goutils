package md

import (
	"fmt"

	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

func init() {
	aa := &oracleActuator{}
	RegisterOQLActuator(aa.GetName(), aa)
}

//公共查询
type oracleActuator struct {
	from   string
	offset int
	limit  int
}

func (oracleActuator) GetName() string {
	return utils.ORM_DRIVER_GODROR
}
func (s *oracleActuator) PlaceholderWrap(dataTypes ...interface{}) string {
	if dataTypes == nil || len(dataTypes) == 0 {
		return "?"
	}
	dType := dataTypes[0].(string)
	if dType == FIELD_TYPE_DATETIME || dType == FIELD_TYPE_DATE {
		if repositories.Default().Dialect().GetName() == utils.ORM_DRIVER_GODROR {
			return fmt.Sprintf("to_date(?,'yyyy-mm-dd')")
		}
	}
	return "?"
}

func (s *oracleActuator) Count(oql *OQL, value interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Pluck(oql *OQL, column string, value interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Take(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Find(oql *OQL, out interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Create(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Update(oql *OQL, data interface{}) *OQL {
	return oql
}
func (s *oracleActuator) Delete(oql *OQL) *OQL {
	return oql
}
