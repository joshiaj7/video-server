package util

var (
	DefaultLimit  = 10
	DefaultOffset = 0
)

type OffsetPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func NewOffsetPagination(limit int, offset int, total int) *OffsetPagination {
	return &OffsetPagination{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
