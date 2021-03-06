package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type Condition struct {
	condition string
	values    []interface{}
}

type ColumnDefinition struct {
	Name         string
	Type		 string
	Nullable     bool
	AutoInc      bool
	ForeignTable string
}

func TableColumn(table Table, col string) string {
	if table.alias == "" {
		return fmt.Sprintf("%s.%s", table.name, col)
	} else {
		return fmt.Sprintf("%s.%s", table.alias, col)
	}

}

func ColInSetCondition(col string, args []interface{}) Condition {
	i := 0

	var argPlaceholders []string
	for i < len(args) {
		argPlaceholders = append(argPlaceholders, "?")
		i += 1
	}

	conditionBody := fmt.Sprintf("%s IN (%s)", col, strings.Join(argPlaceholders, ", "))

	return  Condition{conditionBody, args}
}

func ColEqCondition(col1 string, col2 string) Condition {
	return Condition{condition: fmt.Sprintf("%s = %s", col1, col2), values: []interface{}{}}
}

func SingleValColEqCondition(col string, value interface{}) Condition {
	return SingleValCondition(fmt.Sprintf("%s = ?", col), value)
}

func SingleValCondition(cond string, value interface{}) Condition {
	return Condition{condition: cond, values: []interface{}{value}}
}

func MultipleValCondition(cond string, values []interface{}) Condition {
	return Condition{condition: cond, values: values}
}

type Table struct {
	name  string
	alias string
}

func TableNoAlias(name string) Table {
	return Table{name: name}
}

func TableWithAlias(name string, alias string) Table {
	return Table{name: name, alias: alias}
}

func tableToString(table Table) string {
	if table.alias == "" {
		return table.name
	} else {
		return fmt.Sprintf("%s AS %s", table.name, table.alias)
	}
}
