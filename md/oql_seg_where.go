package md

type OQLWhere interface {
	Where(query string, args ...interface{}) OQLWhere
	OrWhere(query string, args ...interface{}) OQLWhere
	And() OQLWhere
	Or() OQLWhere
	GetArgs() []interface{}
	GetQuery() string
	GetLogical() string
	String() string
}
type oqlWhere struct {
	//字段与操作号之间需要有空格
	//示例1: Org =? ; Org in (?) ;$$Org =?  and ($$Period = ?  or $$Period = ? )
	//示例2：abs($$Qty)>$$TempQty + ?
	Query    string
	Logical  string //and or
	Sequence int
	Children []OQLWhere
	Expr     string
	Args     []interface{}
}

func NewOQLWhere(query string, args ...interface{}) OQLWhere {
	return &oqlWhere{Query: query, Args: args, Logical: "and"}
}
func (m oqlWhere) String() string {
	return m.Query
}
func (m oqlWhere) GetArgs() []interface{} {
	return m.Args
}
func (m oqlWhere) GetQuery() string {
	return m.Query
}
func (m oqlWhere) GetLogical() string {
	return m.Logical
}
func (m *oqlWhere) Where(query string, args ...interface{}) OQLWhere {
	if m.Children == nil {
		m.Children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{Query: query, Args: args, Logical: "and"}
	m.Children = append(m.Children, item)
	return m
}
func (m *oqlWhere) OrWhere(query string, args ...interface{}) OQLWhere {
	if m.Children == nil {
		m.Children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{Query: query, Args: args, Logical: "or"}
	m.Children = append(m.Children, item)
	return m
}
func (m *oqlWhere) And() OQLWhere {
	if m.Children == nil {
		m.Children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{Logical: "and"}
	m.Children = append(m.Children, item)
	return item
}
func (m *oqlWhere) Or() OQLWhere {
	if m.Children == nil {
		m.Children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{Logical: "or"}
	m.Children = append(m.Children, item)
	return item
}
