package repository

import (
	"fmt"
	"strings"
	"yotudo/src/model"
)

const (
	DefaultDateFormat = "2006-01-02_15:04:05"
)

func appendQueryWithFilter(filter string, queryBuilder *strings.Builder, args *[]any) {
	if strings.TrimSpace(filter) != "" {
		queryBuilder.WriteString(" WHERE name LIKE '%' || ? || '%'")
		*args = append(*args, filter)
	}
}

func appendQueryWithSort(sort []model.Sort, queryBuilder *strings.Builder) {
	if len(sort) > 0 {
		for i, s := range sort {
			queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", s.Key, s.DirString()))

			if i != len(sort)-1 {
				queryBuilder.WriteString(",")
			}
		}
	}
}

func appendQueryWithPagination(page *model.Page, queryBuilder *strings.Builder, args *[]any) {
	if page.Size > 0 && page.Page >= 0 {
		queryBuilder.WriteString(" LIMIT ? OFFSET ?")
		*args = append(*args, page.Size, page.Page*page.Size)
	}
}

func inClause(ids []int64, extraArgs ...any) (string, []any) {
	args := make([]any, len(ids)+len(extraArgs))
	ph := make([]string, len(ids))

	if len(extraArgs) > 0 {
		copy(args, extraArgs)
	}

	for i, v := range ids {
		args[i+len(extraArgs)] = v
		ph[i] = "?"
	}

	return strings.Join(ph, ","), args
}
