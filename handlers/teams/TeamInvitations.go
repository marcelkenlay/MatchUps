package teams

import (
	. "../users"
	. "../db"
)

func AddTeamInvitations(teamId int, userIds []int) error {
	columns := []string{TEAM_ID_COL, USER_ID_COL}

	var rows [][]interface{}

	for _, userId := range userIds {
		rows = append(rows, []interface{}{userId})
	}

	return InsertRowsIntoTable(TEAM_INVITATIONS_TABLE_NAME, columns, rows)
}
