package db

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
	"strings"
)

func SelectEntireRowFromTable(table string, conditions []Condition)  (row *sql.Row, err error) {
	return SelectRowFromTable(table, []string{"*"}, conditions)
}

func SelectRowFromTable(table string, columns []string, conditions []Condition)  (row *sql.Row, err error) {
	return BuildSelectFromWhere(columns, []Table{TableNoAlias(table)}, conditions).SelectRow()
}

func SelectRowsFromTable(table string, columns []string, conditions []Condition)  (rows *sql.Rows, err error) {
	return SelectRowsFromTables([]Table{TableNoAlias(table)}, columns, conditions)
}

func SelectRowsFromTables(tables []Table, columns []string, conditions []Condition) (rows *sql.Rows, err error) {
	return BuildSelectFromWhere(columns, tables, conditions).SelectRows()
}





// SELECT QUERY BUILDER

type SelectBuilder sq.SelectBuilder


func BuildSelectFromWhere(columns []string, tables []Table, conditions []Condition) SelectBuilder {
	columnsText := strings.Join(columns, `, `)
	qb := QueryBuilder().Select(columnsText)
	qb = qb.From(tableToString(tables[0]))
	i := 1
	for i < len(tables) {
		qb.Join(tableToString(tables[i]))
	}
	for _, condition := range conditions {
		qb = qb.Where(sq.Expr(condition.condition, condition.values...))
	}
	return SelectBuilder(qb)
}

func (sb SelectBuilder) WithOrdering(orderingCols []string) SelectBuilder {
	ssb := sq.SelectBuilder(sb).OrderBy(orderingCols...)
	return SelectBuilder(ssb)
}

func (sb SelectBuilder) WhereIn(col string, innerQuery SelectBuilder) SelectBuilder {
	psql, args, err := innerQuery.toSql()

	if err != nil {
		log.Println("Error building inner query for WhereIn")
		return SelectBuilder{}
	}

	innerSelect := fmt.Sprintf("%s IN (%s)", col, psql)
	ssb := sq.SelectBuilder(sb).Where(innerSelect, args...)
	return SelectBuilder(ssb)
}

func (sb SelectBuilder) toSql() (psql string, args []interface{}, err error) {
	return  sq.SelectBuilder(sb).ToSql()
}

func (sb SelectBuilder) SelectRow() (row *sql.Row, err error) {
	psql, args, err := sq.SelectBuilder(sb).ToSql()

	if err != nil {
		log.Println("Error building select query", psql, err)
		return nil, err
	}

	return Database.QueryRow(psql, args...), nil
}

func (sb SelectBuilder) SelectRows() (row *sql.Rows, err error) {
	psql, args, err := sq.SelectBuilder(sb).ToSql()

	if err != nil {
		log.Println("Error building select query", psql, err)
		return nil, err
	}

	return Database.Query(psql, args...)
}