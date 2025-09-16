package dto

type Pagination struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

type PaginationMetadata struct {
	Pagination
	Total int `json:"total"`
}
