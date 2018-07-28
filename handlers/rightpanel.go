package handlers

import (
	"net/http"
  _ "github.com/lib/pq"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
  "strconv"
	. "../utils"
)

type UpcomingFixture struct {
	HomeTeam string
	AwayTeam string
	Sport    string
	LocLat   float64
	LocLng   float64
	Date     string
}

type SubmittedScore struct {
	HomeScore   int
	AwayScore   int
	SubmitterID int
}

type PreviousFixture struct {
  HomeTeam  string
  AwayTeam  string
  HomeScore int
  AwayScore int
  Sport     string
}

var GetFixtureDetails = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?fixture_id=x)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	fixtureID := (strings.Split(getquery, "=")[1])

	fields := "t1.team_name, t2.team_name, sport, loc_lat, loc_lng, date"
	tables := "upcoming_fixtures JOIN team_names AS t1 ON home_id=t1.team_id JOIN team_names AS t2 ON away_id=t2.team_id"

	query := fmt.Sprintf("SELECT %s FROM %s WHERE fixture_id=%s;", fields, tables, fixtureID)

	rows, err := Database.Query(query)
	CheckErr(err)

	// Add the only database hit to the result
	if (rows.Next()) {
		data := UpcomingFixture{}
		err = rows.Scan(
		&data.HomeTeam,
		&data.AwayTeam,
		&data.Sport,
		&data.LocLat,
		&data.LocLng,
		&data.Date)

		j,_ := json.Marshal(data) // Convert the list of DB hits to a JSON
		// fmt.Println(">", string(j))
		fmt.Fprintln(writer, string(j)) // Write the result to the sender
	} else {
		fmt.Fprintln(writer, "fail")
	}
})

var GetSubmittedScore = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain fixtureID (query is of the form ?fixture_id=x)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	fixtureID := (strings.Split(getquery, "=")[1])

  // Check if the result has beeen recorded
  query := fmt.Sprintf("SELECT * FROM previous_fixtures WHERE fixture_id=%s;", fixtureID)
  rows, err := Database.Query(query)
  CheckErr(err)

  if (rows.Next()) {
    fmt.Fprintln(writer, "played") // Write the result to the sender
    return
  }

	query = fmt.Sprintf("SELECT submitting_team_id, home_score, away_score FROM submitted_results WHERE fixture_id=%s;", fixtureID)

	rows, err = Database.Query(query)
	CheckErr(err)

	// Add the only database hit to the result
	data := SubmittedScore{}

	if (rows.Next()) {
		err = rows.Scan(
			&data.SubmitterID,
			&data.HomeScore,
			&data.AwayScore)
	} else {
		// No DB hit so set defaults
		data.SubmitterID = -1
		data.HomeScore = -1
		data.AwayScore = -1
	}

	j,_ := json.Marshal(data) // Convert the list of DB hits to a JSON
	// fmt.Println(string(j))
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

var AcceptSubmittedScore = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?fixture_id=x)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	fixtureID := (strings.Split(getquery, "=")[1])

	// Get the score from the submitted results
	query := fmt.Sprintf("SELECT home_score, away_score FROM submitted_results WHERE fixture_id=%s;", fixtureID)
	rows, err := Database.Query(query)
	CheckErr(err)

	var home_score int
	var away_score int

	if (rows.Next()) {
		rows.Scan(&home_score, &away_score)
	} else {
		//Something has gone very wrong
		fmt.Println("chat.go cant find submitted result, ", fixtureID)
		return
	}

	// Remove entry from submitted results
	query = fmt.Sprintf("DELETE FROM submitted_results WHERE fixture_id=%s;", fixtureID)
	_, err = Database.Query(query)
	CheckErr(err)

	// Get details of fixture from upcoming_fixtures
  fields := "home_id, away_id, sport, loc_lat, loc_lng, date"
	query = fmt.Sprintf("SELECT %s FROM upcoming_fixtures WHERE fixture_id=%s;", fields, fixtureID)
	rows, err = Database.Query(query)
	CheckErr(err)

	if (rows.Next()) {
		var home_id int
		var away_id int
		var sport string
		var loc_lat float64
		var loc_lng float64
		var date string

		rows.Scan(
			&home_id,
			&away_id,
			&sport,
			&loc_lat,
			&loc_lng,
			&date)

		value := fmt.Sprintf("%s, %d, %d, '%s', '%s', %d, %d, %f, %f",
			                    fixtureID, home_id, away_id, sport, date, home_score, away_score, loc_lat, loc_lng)
		query = fmt.Sprintf("INSERT INTO previous_fixtures VALUES (%s);", value)
		_, err = Database.Query(query)
		CheckErr(err)

		// Remove entry from upcoming fixtures
		query = fmt.Sprintf("DELETE FROM upcoming_fixtures WHERE fixture_id=%s;", fixtureID)
		_, err = Database.Query(query)
		CheckErr(err)
	} else {
		//Something has gone very wrong
		fmt.Println("chat.go cant find fixture, ", fixtureID)
		return
	}
})

var RejectSubmittedScore = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?fixture_id=x)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	fixtureID := (strings.Split(getquery, "=")[1])

	// Remove entry from submitted results
	query := fmt.Sprintf("DELETE FROM submitted_results WHERE fixture_id=%s;", fixtureID)
	_, err = Database.Query(query)
	CheckErr(err)
})

var SubmitScore = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain info (query is of the form ?fixture_id=x&scH=y&scA=z)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)

  splitAtEq       := strings.Split(getquery, "=")
  fixtureID       := strings.Split(splitAtEq[1], "&")[0]
  homeScoreString := strings.Split(splitAtEq[2], "&")[0]
	awayScoreString := strings.Split(splitAtEq[3], "&")[0]

	query := fmt.Sprintf("INSERT INTO submitted_results VALUES (%s, 200, %s, %s);", fixtureID, homeScoreString, awayScoreString)

	_, err = Database.Query(query)

  if (err != nil) {
    fmt.Fprintln(writer, "fail") // Write the result to the sender
  }
})

var GetPreviousGame = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain fixtureID (query is of the form ?curr=x&opp=y)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)

  splitAtEq       := strings.Split(getquery, "=")
  currTeamString := strings.Split(splitAtEq[1], "&")[0]
	oppoTeamString := splitAtEq[2]

  // Check if the result has beeen recorded
  fields := fmt.Sprintf("home_id, away_id, score_home, score_away, sport")
  cond := fmt.Sprintf("(home_id=%s AND away_id=%s) OR (home_id=%s AND away_id=%s)",
                        currTeamString, oppoTeamString, oppoTeamString, currTeamString)
  query := fmt.Sprintf("SELECT %s FROM previous_fixtures WHERE %s ORDER BY date DESC LIMIT 1;", fields, cond)
  rows, err := Database.Query(query)
  CheckErr(err)

  if (rows.Next()) {
    data := PreviousFixture{}

    rows.Scan(
      &data.HomeTeam,
      &data.AwayTeam,
      &data.HomeScore,
      &data.AwayScore,
      &data.Sport)

    homeID, _ := strconv.ParseInt(data.HomeTeam, 10, 64)
    awayID, _ := strconv.ParseInt(data.AwayTeam, 10, 64)
    data.HomeTeam = GetTeamNameFromTeamID(homeID)
    data.AwayTeam = GetTeamNameFromTeamID(awayID)

    j,_ := json.Marshal(data) // Convert the list of DB hits to a JSON
    // fmt.Println(">>>", string(j))
    fmt.Fprintln(writer, string(j)) // Write the result to the sender
  } else {
    fmt.Fprintln(writer, "fail") // Write the result to the sender
  }
})

func GetTeamNameFromTeamID(teamID int64) (teamname string) {
	// Obtain userID
	query := fmt.Sprintf("SELECT team_name FROM team_names WHERE team_id=%d", teamID)
  row, err := Database.Query(query)
  CheckErr(err)

  if (row.Next()) {
    row.Scan(&teamname)
  } else {
		// Username error
		teamname = "ERR"; // Failure value TODO: make front end handle this
		fmt.Println("Unrecognised teamID (GetTeamNameFromTeamID), ", teamID)
	}

	return teamname
}
