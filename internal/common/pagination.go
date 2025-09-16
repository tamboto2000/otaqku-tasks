package common

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
