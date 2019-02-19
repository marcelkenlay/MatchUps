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
	LAT string `json:"lat"`
	LNG string `json:"lng"`
}

var TEAMS_TABLE = "teams"
var TEAM_MEMBERS_TABLE = "team_members"
var TEAM_ID = "team_id"
var TEAM_NAME = "name"
var TM_USER_ID  = "user_id"

func GetUsersTeams(userSession UserSessionCookie) {

	userId := GetUserIdFromSession(userSession)

	// Tables for query
	teamMembersTable := TableNoAlias(TEAM_MEMBERS_TABLE)
	teamsTable := TableNoAlias(TEAMS_TABLE)
	usersTable := TableNoAlias(USERS_TABLE)


	innerSelectCols := []string{TEAM_ID}
	innerSelectConds := []Condition{SingleValColEqCondition(USER_ID, userId)}
	innerSelectTables := []Table{TableNoAlias(TEAM_MEMBERS_TABLE)}
	usersTeams := BuildSelectFromWhere(innerSelectCols, innerSelectTables, innerSelectConds)


	// Data For Main Query
	selectCols := []string{TableColumn(teamsTable, TEAM_ID), TableColumn(teamsTable, TEAM_NAME),
		TableColumn(usersTable, USERNAME), TableColumn(usersTable, NAME), TableColumn(usersTable, LOC_LAT),
		TableColumn(usersTable, LOC_LNG)}

	selectTables := []Table{teamsTable, teamMembersTable, usersTable}

	selectConds := []Condition{ColEqCondition(TableColumn(teamsTable, TEAM_ID), TableColumn(teamMembersTable, TEAM_ID)),
		ColEqCondition(TableColumn(teamMembersTable, TM_USER_ID), TableColumn(usersTable, USER_ID))}

	// Build the select clause and then select the rows
	rows, err := BuildSelectFromWhere(selectCols, selectTables, selectConds).
			WhereIn(TEAM_ID, usersTeams).WithOrdering([]string{TEAM_ID}).SelectRows()

	if err != nil {
		log.Println("Error selecting rows for get user teams")
	}

	var teams []Team
	var curTeam Team
	
	for rows.Next() {
		var teamId int
		var teamName string
		player := Player{LOCATION: Location{} }

		_ = rows.Scan(&teamId, &teamName, &player.USERNAME, &player.NAME, &player.LOCATION.LAT, &player.LOCATION.LNG)

		if curTeam.ID != teamId {
			teams = append(teams, curTeam)
			curTeam = Team{ID:teamId, NAME:teamName}
		}

		curTeam.PLAYERS = append(curTeam.PLAYERS, player)
	}

}