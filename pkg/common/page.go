package common

import "math"

func CalculateTotalPages(totalItems uint, pageSize uint) uint {
	return uint(math.Ceil(float64(totalItems) / float64(pageSize)))
}
