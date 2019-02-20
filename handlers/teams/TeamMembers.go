package teams

import (
	. "../db"
	. "../users"
	"log"
)

type Team struct {
	ID      int      `json:"id"`
	NAME    string   `json:"name"`
	IMAGE   string   `json:"image"`
	PLAYERS []Player `json:"players"`
}

type Player struct {
	USERNAME string   `json:"username"`
	NAME     string   `json:"name"`
	IMAGE    string   `json:"image"`
	LOCATION Location `json:"location"`
}

type Location struct {
	LAT float64 `json:"lat"`
	LNG float64 `json:"lng"`
}


func AddMemberToTeam(teamId int, userId int) error {
	columns := []string{USER_ID_COL, TEAM_ID_COL}
	args := []interface{}{userId, teamId}
	return InsertRowIntoTable(TEAM_MEMBERS_TABLE_NAME, columns, args)
}


func GetUsersTeams(userSession UserSessionCookie) []Team {

	userId := GetUserIdFromSession(userSession)

	// Tables for query
	teamMembersTable := TableNoAlias(TEAM_MEMBERS_TABLE_NAME)
	teamsTable := TableNoAlias(TEAMS_TABLE_NAME)
	usersTable := TableNoAlias(USERS_TABLE)

	innerSelectCols := []string{TEAM_ID_COL}
	innerSelectConds := []Condition{SingleValColEqCondition(USER_ID_COL, userId)}
	innerSelectTables := []Table{TableNoAlias(TEAM_MEMBERS_TABLE_NAME)}
	usersTeams := BuildSelectFromWhere(innerSelectCols, innerSelectTables, innerSelectConds)

	// Data For Main Query
	selectCols := []string{TableColumn(teamsTable, TEAM_ID_COL), TableColumn(teamsTable, TEAM_NAME_COL),
		TableColumn(usersTable, USERNAME_COL), TableColumn(usersTable, NAME_COL), TableColumn(usersTable, LOC_LAT_COL),
		TableColumn(usersTable, LOC_LNG_COL)}

	selectTables := []Table{teamsTable, teamMembersTable, usersTable}

	selectConds := []Condition{
			ColEqCondition(TableColumn(teamsTable, TEAM_ID_COL), TableColumn(teamMembersTable, TEAM_ID_COL)),
			ColEqCondition(TableColumn(teamMembersTable, USER_ID_COL), TableColumn(usersTable, USER_ID_COL) )}

	// Build the select clause and then select the rows
	rows, err := BuildSelectFromWhere(selectCols, selectTables, selectConds).
		WhereIn(TEAM_ID_COL, usersTeams).WithOrdering([]string{TEAM_ID_COL}).SelectRows()

	if err != nil {
		log.Println("Error selecting rows for get user teams")
	}

	var teams []Team
	var curTeam Team

	for rows.Next() {
		var teamId int
		var teamName string
		player := Player{LOCATION: Location{}}

		_ = rows.Scan(&teamId, &teamName, &player.USERNAME, &player.NAME, &player.LOCATION.LAT, &player.LOCATION.LNG)

		if curTeam.ID != teamId {
			teams = append(teams, curTeam)
			curTeam = Team{ID: teamId, NAME: teamName}
		}

		curTeam.PLAYERS = append(curTeam.PLAYERS, player)
	}

	return teams
}
