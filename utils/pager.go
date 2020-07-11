package utils

type Pager struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	LastPage int `json:"last_page"`
	Items    int `json:"items"` //当前条数
	Total    int `json:"total"` //总记录数
}
