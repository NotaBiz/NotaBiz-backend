package helper

import (
	"log"
	"net/url"
	"strconv"
)

func GetPaginationParam(query string) (page int, pageSize int, err error) {
	queryParams, err := url.ParseQuery(query)
	if err != nil {
		log.Println("Error parsing query string:", err)
		return 0, 0, err
	}
	pageFilter, ok := queryParams["page"]
	if ok {
		page, err = strconv.Atoi(pageFilter[0])
		if err != nil {
			return 0, 0, err
		}
	}
	sizeFilter, ok := queryParams["page_size"]
	if ok {
		pageSize, err = strconv.Atoi(sizeFilter[0])
		if err != nil {
			return 0, 0, err
		}
	}
	if pageSize <= 0 || page <= 0 {
		pageSize = 999999999999999
		page = 1
	}
	page = (page - 1) * pageSize
	return
}

func GetTotalPage(totalRows string, rowsPerPage int) int {
	var totalPage int
	tr, _ := strconv.Atoi(totalRows)
	rest := tr % rowsPerPage
	totalPage = tr / rowsPerPage
	if rest != 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

func GetHasNext(currentPage, totalRows, rowsPerPage int) bool {
    if rowsPerPage <= 0 {
        return false
    }
    return currentPage*rowsPerPage < totalRows
}

func GetHasPrev(currentPage int) bool {
	return currentPage > 1
}