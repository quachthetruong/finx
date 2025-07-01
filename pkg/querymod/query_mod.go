package querymod

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
)

type InConstraint interface {
	int | int64 | uint64 | string
}

func In[T InConstraint](ss []T) []postgres.Expression {
	if len(ss) == 0 {
		return nil
	}
	ee := make([]postgres.Expression, 0, len(ss))
	switch any(ss[0]).(type) {
	case int:
		for _, s := range ss {
			ee = append(ee, postgres.Int64(int64(any(s).(int))))
		}
	case int64:
		for _, s := range ss {
			ee = append(ee, postgres.Int64(any(s).(int64)))
		}
	case uint64:
		for _, s := range ss {
			ee = append(ee, postgres.Uint64(any(s).(uint64)))
		}
	case string:
		for _, s := range ss {
			ee = append(ee, postgres.String(any(s).(string)))
		}
	}
	return ee
}

type GetQm struct {
	ForUpdate bool
}

type GetOption func(g *GetQm)

func WithLock() GetOption {
	return func(g *GetQm) {
		g.ForUpdate = true
	}
}

func RawSubQuery(stm postgres.Statement) postgres.Expression {
	return postgres.Raw(strings.TrimSuffix(strings.TrimSpace(ReplacePlaceholders(stm.Sql())), ";"))
}

func ReplacePlaceholders(query string, values []interface{}) string {
	for i, v := range values {
		query = strings.Replace(query, fmt.Sprintf("$%d", i+1), formatPlaceholder(v), 1)
	}
	return query
}

func formatPlaceholder(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	// Handle other data types as needed
	default:
		return fmt.Sprintf("%v", v)
	}
}
