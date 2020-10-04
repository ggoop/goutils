package md

type OQLOrder interface {
	GetQuery() string
	GetSequence() int
	GetArgs() []interface{}
}

type oqlOrder struct {
	Origin   interface{}
	Query    string
	Sequence int
	Expr     string
	Args     []interface{}
}

func (m oqlOrder) GetQuery() string {
	return m.Query
}
func (m oqlOrder) GetSequence() int {
	return m.Sequence
}
func (m oqlOrder) GetArgs() []interface{} {
	return m.Args
}
