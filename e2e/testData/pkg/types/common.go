package types

type JustID struct {
	ID int `json:"id,omitempty" position:"path"`
}

type PagerReq struct {
	PageSize int `json:"page_size,omitempty" position:"query"`
	Page     int `json:"page,omitempty" position:"query"`
}

type PagerRes struct {
	TotalCount int         `json:"total_count,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}
