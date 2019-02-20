package db

import (
	"fmt"
	"strings"
)

func CreateTable(name string, columns []ColumnDefinition, primaryColumns []ColumnDefinition) error {
	var columnDefs []string
	for _, column := range columns {
		var columnDef []string

		columnDef = append(columnDef, column.Name)

		if column.AutoInc {
			columnDef = append(columnDef, "serial")
		} else {
			columnDef = append(columnDef, column.Type)
		}

		if !column.Nullable {
			columnDef = append(columnDef, "not null")
		}

		if column.ForeignTable != "" {
			fkConstraint := fmt.Sprintf("constraint %s_%s_%s_fk references %s",
				name, column.ForeignTable, column.Name, column.ForeignTable)
			columnDef = append(columnDef, fkConstraint)
		}

		columnDefs = append(columnDefs, strings.Join(columnDef, " "))
	}

	if len(primaryColumns) > 0 {
		var primaryColumnNames []string
		for _, primaryColumn := range primaryColumns {
			primaryColumnNames = append(primaryColumnNames, primaryColumn.Name)
		}

		fkConstraint := fmt.Sprintf("constraint %s_fk unique (%s)",
			name, strings.Join(primaryColumnNames, ","))

		columnDefs = append(columnDefs, fkConstraint)
	}


	query := fmt.Sprintf("create table %s ( %s )", name, strings.Join(columnDefs, ", "))

	_, err := Database.Query(query)

	return err
}