package helper

import (
	"math"
)

func Pagination[T comparable](x []T, page int, perPage int) (data []T, currentPage int, lastPage int) {
	lastPage = int(math.Ceil(float64(len(x)) / float64(perPage)))
	currentPage = page

	if page < 1 {
		page = 1
	} else if lastPage < page {
		page = lastPage
		currentPage = page
		return []T{}, currentPage, lastPage
	}

	if page == lastPage {
		data = x[(page-1)*perPage:]
	} else {
		data = x[(page-1)*perPage : page*perPage]
	}

	return
}
