package handlers

import (
	. "./chats"
	. "./db"
	. "./teams"
	. "./users"
	. "./utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type NewTeamInfo struct {
	TeamName         string
	CaptainID        int
	InviteeUsernames []string
}

var AddNewTeam = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var teamInfo NewTeamInfo
	err := decoder.Decode(&teamInfo)

	if err != nil {
		panic(err)
	}

	userSessionCookie := BuildUserSessionFromRequest(request)

	if !IsTeamNameAvailable(teamInfo.TeamName) {
		_, _ = fmt.Fprintln(writer, -1)
		return
	}

	teamId, err := AddNewTeamAndRetrieveId(teamInfo.TeamName)

	if err != nil {
		log.Println("Error adding team: ", teamId)
		_, _ = fmt.Fprintln(writer, -1)
		return
	}

	userId := GetUserIdFromSession(userSessionCookie)

	_ = AddMemberToTeam(teamId, userId)

	_ = AddTeamCaptain(teamId, userId)

	var location Location
	location.LAT, location.LNG = GetUserLatLngFromId(userId)

	_ = AddTeamLocation(teamId, location)

	_ = AddTeamAvailability(teamId)

	_ = CreateTeamChat(teamId)

	_, _ = fmt.Fprintln(writer, teamId) // Write team id to signal success to the sender
})

var GetTeams = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userSessionCookie := BuildUserSessionFromRequest(r)

	teams := GetUsersTeams(userSessionCookie)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(w).Encode(teams)
	CheckErr(err)
})

var GetInvitations = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//get username
	teams := []Team{}
	username := mux.Vars(r)["username"]
	query := fmt.Sprintf("SELECT team_name "+
		"FROM users "+
		"JOIN team_invitations on users.user_id=team_invitations.player_id "+
		"JOIN team_names on team_names.team_id=team_invitations.team_id "+
		"WHERE username='%s'", username)
	rows, err := Database.Query(query)
	CheckErr(err)

	for rows.Next() {
		team := Team{}
		err := rows.Scan(&team.NAME)
		CheckErr(err)
		query = fmt.Sprintf("SELECT username,users.loc_lat,users.loc_lng "+
			"FROM team_members "+
			"JOIN team_names on team_members.team_id = team_names.team_id "+
			"JOIN users on team_members.user_id=users.user_id "+
			"WHERE team_name='%s';", team.NAME)
		users, err := Database.Query(query)
		players := []Player{}
		//for each player retrieve location
		for users.Next() {
			player := Player{}
			loc := Location{}
			err := users.Scan(&player.NAME, &loc.LAT, &loc.LNG)
			CheckErr(err)
			player.LOCATION = loc
			players = append(players, player)
		}
		team.PLAYERS = players
		teams = append(teams, team)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(teams)
})

type QueryMatch struct {
	UserID   int
	Username string
	FullName string
}

var GetUsernameMatches = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// Obtain pattern to match (query is of the form ?pattern=)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	pattern := strings.Split(getquery, "=")[1]

	query := fmt.Sprintf("SELECT user_id, username, name FROM users WHERE UPPER(username) LIKE '%s%s';", strings.ToUpper(pattern), "%")
	// fmt.Println(query)
	rows, err := Database.Query(query)
	CheckErr(err)

	var result []QueryMatch
	for rows.Next() {
		data := QueryMatch{}
		err = rows.Scan(&data.UserID, &data.Username, &data.FullName)
		CheckErr(err)
		result = append(result, data)
	}

	j, _ := json.Marshal(result)    // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

type TeamInvInfo struct {
	TeamID   int
	Invitees []int
}

var SendInvitations = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var teamInvInfo TeamInvInfo
	err := decoder.Decode(&teamInvInfo)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

	//Add invitations
	for _, invitee := range teamInvInfo.Invitees {
		query := fmt.Sprintf("INSERT INTO team_invitations VALUES(%d, %d);",
			teamInvInfo.TeamID, invitee)
		// fmt.Println(query)
		_, err = Database.Query(query)
		CheckErr(err)
	}

})

var AddPlayerToTeam = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	//get user-id and team-id
	var userId, teamId int
	query := fmt.Sprintf("select user_id "+
		"FROM users where username='%s'", vars["username"])
	err := Database.QueryRow(query).Scan(&userId)
	CheckErr(err)
	query = fmt.Sprintf("select team_id "+
		"FROM team_names where team_name='%s'", vars["teamname"])
	err = Database.QueryRow(query).Scan(&teamId)
	CheckErr(err)
	//insert into team_members
	query = fmt.Sprintf("INSERT INTO team_members VALUES('%d', '%d');",
		teamId, userId)
	_, err = Database.Query(query)
	CheckErr(err)

	//remove team from team_invitations
	query = fmt.Sprintf(
		"DELETE FROM team_invitations "+
			"WHERE team_id=%d AND player_id=%d",
		teamId, userId)
	_, err = Database.Query(query)
	CheckErr(err)

	//send updated team
	team := Team{}
	query = fmt.Sprintf("SELECT username,users.loc_lat,users.loc_lng "+
		"FROM team_members "+
		"JOIN team_names on team_members.team_id = team_names.team_id "+
		"JOIN users on team_members.user_id=users.user_id "+
		"WHERE team_name='%s';", vars["teamname"])
	users, err := Database.Query(query)
	players := []Player{}
	//for each player retrieve location
	for users.Next() {
		player := Player{}
		loc := Location{}
		err := users.Scan(&player.NAME, &loc.LAT, &loc.LNG)
		CheckErr(err)
		player.LOCATION = loc
		players = append(players, player)
	}
	team.PLAYERS = players
	team.NAME = vars["teamname"]
	json.NewEncoder(writer).Encode(team)

	// Recalculate the team's location and availability
	RecalculateTeamAvailability(teamId)
	RecalculateTeamLocation(teamId)
})

var DeleteInvitation = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	//get user-id and team-id
	var userId, teamId int
	query := fmt.Sprintf("select user_id "+
		"FROM users where username='%s'", vars["username"])
	err := Database.QueryRow(query).Scan(&userId)
	CheckErr(err)
	query = fmt.Sprintf("select team_id "+
		"FROM team_names where team_name='%s'", vars["teamname"])
	err = Database.QueryRow(query).Scan(&teamId)
	CheckErr(err)

	//remove team from team_invitations
	query = fmt.Sprintf(
		"DELETE FROM team_invitations "+
			"WHERE team_id=%d AND player_id=%d",
		teamId, userId)
	_, err = Database.Query(query)
	CheckErr(err)
	writer.WriteHeader(http.StatusOK)
})

var GetTeamNames = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	team_id := strings.Split(getquery, "=")[1]

	// Run query
	query := fmt.Sprintf("SELECT users.name FROM team_members "+
		"NATURAL INNER JOIN users "+
		" WHERE team_members.team_id=%s;", team_id)
	rows, err := Database.Query(query)
	CheckErr(err)

	var result []string
	// Add every database hit to the result
	for rows.Next() {
		var member string
		err = rows.Scan(&member)
		result = append(result, member)
	}

	j, _ := json.Marshal(result)    // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})
