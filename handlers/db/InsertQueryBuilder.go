package db

import (
	"errors"
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

	qb, err := buildInsertQuery(table, columns, vals)

	if err != nil {
		return nil, err
	}

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

	qb, err := buildInsertQuery(table, columns, vals)

	if err != nil {
		return err
	}

	psql, args, err := qb.ToSql()

	_, err = Database.Query(psql, args...)

	return err
}

func buildInsertQuery(table string, columns []string, vals []interface{}) (qb sq.InsertBuilder, err error) {
	qb = QueryBuilder().Insert(table)

	if len(columns) != len(vals) {
		print("Columns and values lengths do not match")
		err = errors.New("columns and values lengths do not match")
		return
	}

	qb = qb.Columns(columns...).Values(vals...)

	return
}
