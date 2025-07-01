package querymod

import (
	"github.com/go-jet/jet/v2/postgres"
)

func ArrayAgg(expr postgres.Expression) postgres.Expression {
	return postgres.Func("array_agg", expr)
}

func ArrayAny(expr postgres.Expression) postgres.Expression {
	return postgres.Func("any", expr)
}

func ArrayAll(expr postgres.Expression) postgres.Expression {
	return postgres.Func("all", expr)
}
