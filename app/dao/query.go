package dao

import "fmt"

// BuildRangeQuery returns query string and query exists or not
func BuildRangeQuery(column string, max, since, invalid int64) (string, bool) {
	var where, maxQuery, sinceQuery, and string
	var query string

	if max != invalid || since != invalid {
		where = "WHERE"
		if max != invalid {
			maxQuery = fmt.Sprintf("%s <= %d", column, max)
		}
		if since != invalid {
			sinceQuery = fmt.Sprintf("%s >= %d", column, since)
		}
		if max != invalid && since != invalid {
			and = "AND"
		}
		query = fmt.Sprintf("%s %s %s %s", where, maxQuery, and, sinceQuery)
	}

	exist := (query != "")

	return query, exist
}
