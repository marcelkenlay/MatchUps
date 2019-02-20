package db

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func InsertRowIntoTableAndRetreiveVal(table string, columns []string, vals []interface{}, retrieve string) (string, error) {

	readCols, err := InsertRowIntoTableAndRetreiveVals(table, columns, vals, []string{retrieve})
	return readCols[0], err
}


func InsertRowIntoTableAndRetreiveVals(table string, columns []string, vals []interface{}, retrieve []string) ([]string, error) {

	qb := buildInsertRowQuery(table, columns, vals)

	qb = qb.Suffix(fmt.Sprintf("RETURNING \"%s\"", strings.Join(retrieve, " , ")))

	psql, args, err := qb.ToSql()

	if err != nil {
		log.Println("Error building INSERT query", psql, err)
		return nil, err
	}

	row := Database.QueryRow(psql, args...)

	readCols := make([]interface{}, len(retrieve))
	writeCols := make([]string, len(retrieve))
	for i := range writeCols {
		readCols[i] = &writeCols[i]
	}

	err = row.Scan(readCols...)

	if err != nil {
		log.Println("Error scanning result", psql, err)
	}

	return writeCols, err

}

func InsertRowIntoTable(table string, columns []string, vals []interface{}) error {
	psql, args, err := buildInsertRowQuery(table, columns, vals).ToSql()

	if err != nil {
		return err
	}

	_, err = Database.Query(psql, args...)

	return err
}

func InsertRowsIntoTable(table string, columns []string, vals [][]interface{}) error {
	psql, args, err :=  buildInsertRowsQuery(table, columns, vals).ToSql()

	if err != nil {
		return err
	}

	_, err = Database.Query(psql, args...)

	return err
}

func buildInsertRowQuery(table string, columns []string, vals []interface{}) sq.InsertBuilder {
	return buildInsertRowsQuery(table, columns, [][]interface{}{vals})
}


func buildInsertRowsQuery(table string, columns []string, vals [][]interface{}) sq.InsertBuilder {
	qb := QueryBuilder().Insert(table)

	qb = qb.Columns(columns...)

	for _, val := range vals {
		qb = qb.Values(val...)
	}

	return qb
}
