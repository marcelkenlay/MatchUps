package teams

import (
	. "../db"
	. "../users"
	"fmt"
	"strconv"
	"strings"
)

func IsTeamNameAvailable(teamName string) bool {

	teamName = strings.ToUpper(teamName)
	conditions := []Condition{SingleValCondition(fmt.Sprintf("UPPER(%s) = ?", TEAM_NAME_COL), teamName)}

	num, err := CountMatchingRows(TEAMS_TABLE_NAME, conditions)

	if err != nil || num > 0 {
		return false
	}

	return true
}

func AddNewTeamAndRetrieveId(teamName string) (int, error) {
	columns := []string{TEAM_NAME_COL}
	vals := []interface{}{teamName}

	id, err := InsertRowIntoTableAndRetreiveVal(TEAMS_TABLE_NAME, columns, vals, TEAM_ID_COL)

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(id)
}

func AddTeamCaptain(teamId int, captainId int) error {
	columns := []string{TEAM_ID_COL, USER_ID_COL}
	vals := []interface{}{teamId, captainId}

	return InsertRowIntoTable(TEAM_CAPTAINS_TABLE_NAME, columns, vals)
}

func AddTeamLocation(teamId int, location Location) error {
	columns := []string{TEAM_ID_COL, LOC_LAT_COL, LOC_LNG_COL}
	vals := []interface{}{teamId, location.LAT, location.LNG}
	return InsertRowIntoTable(TEAM_LOCATIONS_TABLE_NAME, columns, vals)
}

func AddTeamAvailability(teamId int) error {
	columns := []string{TEAM_ID_COL}
	vals := []interface{}{teamId}
	return InsertRowIntoTable(TEAM_AVAILABILITY_TABLE, columns, vals)
}