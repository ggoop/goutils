package md

type oqlEntity struct {
	Path     string
	Entity   *MDEntity
	Sequence int
	IsMain   bool
	Alias    string
}
type oqlField struct {
	Entity *oqlEntity
	Field  *MDField
	Path   string
}

type OQLFrom struct {
	Query string
	Alias string
	Args  []interface{}
	expr  string
}

type OQLJoin struct {
	Type      OQLJoinType
	Query     string
	Alias     string
	Condition string
	Args      []interface{}
	expr      string
}
type OQLSelect struct {
	Query string
	Alias string
	Args  []interface{}
	expr  string
}
type OQLGroup struct {
	Query string
	Args  []interface{}
	expr  string
}
type OQLOrder struct {
	Query string
	Order OQLOrderType
	Args  []interface{}
	expr  string
}
