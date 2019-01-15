package query

type IWhere interface {
	Where(query string, args ...interface{}) IWhere
	OrWhere(query string, args ...interface{}) IWhere
	And() IWhere
	Or() IWhere
	GetArgs() []interface{}
	GetQuery() string
	GetLogical() string
}
type oqlWhere struct {
	Query    string
	Logical  string //and or
	Sequence int
	Children []*oqlWhere
	Expr     string
	Args     []interface{}
}

func (m *oqlWhere) GetArgs() []interface{} {
	return m.Args
}
func (m *oqlWhere) GetQuery() string {
	return m.Query
}
func (m *oqlWhere) GetLogical() string {
	return m.Logical
}
func (m *oqlWhere) Where(query string, args ...interface{}) IWhere {
	if m.Children == nil {
		m.Children = make([]*oqlWhere, 0)
	}
	item := &oqlWhere{Query: query, Args: args, Logical: "and"}
	m.Children = append(m.Children, item)
	return m
}
func (m *oqlWhere) OrWhere(query string, args ...interface{}) IWhere {
	if m.Children == nil {
		m.Children = make([]*oqlWhere, 0)
	}
	item := &oqlWhere{Query: query, Args: args, Logical: "or"}
	m.Children = append(m.Children, item)
	return m
}
func (m *oqlWhere) And() IWhere {
	if m.Children == nil {
		m.Children = make([]*oqlWhere, 0)
	}
	item := &oqlWhere{Logical: "and"}
	m.Children = append(m.Children, item)
	return item
}
func (m *oqlWhere) Or() IWhere {
	if m.Children == nil {
		m.Children = make([]*oqlWhere, 0)
	}
	item := &oqlWhere{Logical: "or"}
	m.Children = append(m.Children, item)
	return item
}
