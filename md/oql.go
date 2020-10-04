package md

const (
	regexp_oql_from   = "([\\S]+)(?i:(?:as|[\\s])+)([\\S]+)|([\\S]+)"
	regexp_oql_select = "([\\S]+.*\\S)(?i:\\s+as+\\s)([\\S]+)|([\\S]+.*[\\S]+)"
	regexp_oql_order  = "(?i)([\\S]+.*\\S)(?:\\s)(desc|asc)|([\\S]+.*[\\S]+)"
)

func GetOQL(names ...string) *OQL {
	return &OQL{}
}
