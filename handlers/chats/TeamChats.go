package chats

import (
	. "../users"
	. "../db"
	"fmt"
)

func CreateTeamChat(teamId int) error {

	table_name := fmt.Sprintf("chats.team%d", teamId)

	sender_col := ColumnDefinition{Name:"user_id", Type: "integer", ForeignTable: USERS_TABLE}
	message_col := ColumnDefinition{Name:"message", Type: "varchar(200)"}
	time_sent_col := ColumnDefinition{Name:"time_sent", Type: "timestamp without time zone"}

	columns := []ColumnDefinition{sender_col, message_col, time_sent_col}
	var primaryColumns []ColumnDefinition

	return CreateTable(table_name, columns, primaryColumns)
}
