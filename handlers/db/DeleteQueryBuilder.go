package db

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

func DeleteRowsFromTable(table string, conditions []Condition)  (rows *sql.Rows, err error) {
	psql, args, err := buildDeleteQuery([]Table{TableNoAlias(table)}, conditions)

	if err != nil {
		return
	}

	return Database.Query(psql, args...)
}

func buildDeleteQuery(tables []Table, conditions []Condition) (string, []interface{}, error) {
	qb := QueryBuilder().Delete(tableToString(tables[0]))

	i := 1
	for i < len(tables) {
		qb.Suffix(fmt.Sprintf("USING %s", tableToString(tables[i])), []interface{}{})
	}

	for _, condition := range conditions {
		qb = qb.Where(sq.Expr(condition.condition, condition.values...))
	}

	psql, args, err := qb.ToSql()

	if err != nil {
		log.Println("Error building select query", psql, err)
		return "", nil, err
	}

	return psql, args, err
}
