package repository

import (
	"strings"
)

const (
	DefaultDateFormat = "2006-01-02_15:04:05"
)

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
