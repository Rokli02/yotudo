package builders

import (
	"fmt"
	"strings"
	"yotudo/src/model"
)

// Query Builder States
const (
	whereConditionAdded int = 1 << iota
	sortAdded
	paginationAdded
	ignoreEndingSemicolon
)

type QueryBuilder struct {
	builder strings.Builder
	args    *[]any
	state   int
}

func NewQueryBuilder(base string, args *[]any) *QueryBuilder {
	qb := &QueryBuilder{
		args: args,
	}

	qb.builder.WriteString(base)

	return qb
}

func (qb *QueryBuilder) Build() string {
	if qb.state&ignoreEndingSemicolon == 0 {
		qb.builder.WriteString(";")
	}

	return qb.builder.String()
}

func (qb *QueryBuilder) WithoutSemicolon() QueryBuilderFinalStage {
	qb.state |= ignoreEndingSemicolon

	return qb
}

func (qb *QueryBuilder) WithFilter(key string, value string) QueryBuilderStage1 {
	if strings.TrimSpace(value) != "" {
		if qb.state&whereConditionAdded != 0 {
			qb.builder.WriteString(fmt.Sprintf(" AND %s LIKE '%%' || ? || '%%'", key))
		} else {
			qb.builder.WriteString(fmt.Sprintf(" WHERE %s LIKE '%%' || ? || '%%'", key))
			qb.state |= whereConditionAdded
		}

		*qb.args = append(*qb.args, value)
	}

	return qb
}

func (qb *QueryBuilder) WithCondition(key string, value any, isValidFuncs ...func(value any) bool) QueryBuilderStage1 {
	for _, isValid := range isValidFuncs {
		if !isValid(value) {
			return qb
		}
	}

	if qb.state&whereConditionAdded != 0 {
		qb.builder.WriteString(fmt.Sprintf(" AND %s=?", key))
	} else {
		qb.builder.WriteString(fmt.Sprintf(" WHERE %s=?", key))
		qb.state |= whereConditionAdded
	}

	*qb.args = append(*qb.args, value)

	return qb
}

func (qb *QueryBuilder) WithSort(sort []model.Sort) QueryBuilderStage3 {
	if len(sort) > 0 {
		for i, s := range sort {
			qb.builder.WriteString(fmt.Sprintf(" ORDER BY %s %s", s.Key, s.DirString()))

			if i != len(sort)-1 {
				qb.builder.WriteString(",")
			}
		}

		qb.state |= sortAdded
	}

	return qb
}

func (qb *QueryBuilder) WithPagination(page *model.Page) QueryBuilderFinalStage {
	if page.Size > 0 && page.Page >= 0 {

		qb.builder.WriteString(" LIMIT ? OFFSET ?")
		*qb.args = append(*qb.args, page.Size, page.Page*page.Size)

		qb.state |= paginationAdded
	}

	return qb
}

type QueryBuilderStage1 interface {
	WithFilter(key string, value string) QueryBuilderStage1
	WithCondition(key string, value any, isValidFuncs ...func(value any) bool) QueryBuilderStage1
	QueryBuilderStage2
}

type QueryBuilderStage2 interface {
	WithSort(sort []model.Sort) QueryBuilderStage3
	QueryBuilderStage3
}

type QueryBuilderStage3 interface {
	WithPagination(page *model.Page) QueryBuilderFinalStage
	QueryBuilderFinalStage
}

type QueryBuilderFinalStage interface {
	WithoutSemicolon() QueryBuilderFinalStage
	Build() string
}
