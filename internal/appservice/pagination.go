package appservice

import "github.com/songwei.ma/talus-mofish/internal/store"

func normalizePageParams(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = store.DefaultPageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
