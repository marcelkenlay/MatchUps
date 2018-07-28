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

type UserInfo struct {
	LocLat  float64
	LocLng  float64
}

type Availability struct {
	Mon  int64
	Tues int64
	Wed  int64
	Thurs int64
	Fri  int64
	Sat  int64
	Sun  int64
}

type Fixture struct {
	Opposition    string
	ForTeam       string
	Sport         string
	LocLat        string
	LocLng        string
	Date          string
	ScoreHome     int
	ScoreAway     int
	IsHome        bool
}

// Query the database for a userID corresponding to a username
func GetUserIDFromUsername(username string) (userID int) {
	// Obtain userID
	query := fmt.Sprintf("SELECT user_id FROM users WHERE username='%s'", username)
  row, err := Database.Query(query)
  CheckErr(err)

  if (row.Next()) {
    row.Scan(&userID)
  } else {
		// Username error
		userID = -1; // Failure value TODO: make front end handle this
		fmt.Println("Unrecognised username (GetUserUpcoming), ", username)
	}

	return userID
}

// Get the user information for the user specified in the url
var GetUserInfo = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := (strings.Split(getquery, "=")[1])

	// Run query
  query := fmt.Sprintf("SELECT loc_lat, loc_lng FROM users WHERE username='%s';", username)
  rows, _ := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	if (rows.Next()) {
		data := UserInfo{}
		err = rows.Scan(&data.LocLat, &data.LocLng)

		j,_ := json.Marshal(data) // Convert the list of DB hits to a JSON
		// fmt.Println(string(j))
		fmt.Fprintln(writer, string(j)) // Write the result to the sender
	} else {
		fmt.Println("No user info DB hit for ", username)
	}
})


// Get the fixtures for the user specified in the url
var GetUserUpcoming = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username=name)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := strings.Split(getquery, "=")[1]

	// Obtain userID
  userID := GetUserIDFromUsername(username)

  // Build queries
	ordering := "ORDER BY date ASC"
	commonQueryFields :=  "sport, loc_lat, loc_lng, date"
  tables := "upcoming_fixtures JOIN team_members"

	homeFields := fmt.Sprintf("away_id, home_id, %s", commonQueryFields)
	homeTableJoinCond := "home_id=team_id"

	awayFields := fmt.Sprintf("home_id, away_id, %s", commonQueryFields)
	awayTableJoinCond := "away_id=team_id"

	// Common text used in initialising JSON responses
	var jsonText = []byte(`[]`)

	homeQuery := fmt.Sprintf("SELECT %s FROM %s ON %s WHERE user_id=%d %s",
		                        homeFields, tables, homeTableJoinCond, userID, ordering)
	awayQuery := fmt.Sprintf("SELECT %s FROM %s ON %s WHERE user_id=%d %s",
		                        awayFields, tables, awayTableJoinCond, userID, ordering)

  // Run the query for home games
	rows, err := Database.Query(homeQuery)
	CheckErr(err)

	// Initialise the json response for all home games
	var teamHome []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamHome)

	// Add every database hit to the result
	for rows.Next() {
		data := Fixture{}
		err = rows.Scan(
			&data.Opposition,
			&data.ForTeam,
			&data.Sport,
			&data.LocLat,
			&data.LocLng,
			&data.Date)

		oppID, _ := strconv.ParseInt(data.Opposition, 10, 64)
		forID, _ := strconv.ParseInt(data.ForTeam, 10, 64)
		data.Opposition = GetTeamNameFromTeamID(oppID)
		data.ForTeam = GetTeamNameFromTeamID(forID)

		data.IsHome = true

		teamHome = append(teamHome, data)
	}

	rows, err = Database.Query(awayQuery)
	CheckErr(err)

	// Initialise the json response for all away games
	var teamAway []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamAway)

	// Add every database hit to the result
	for rows.Next() {
		data := Fixture{}
		err = rows.Scan(
			&data.Opposition,
			&data.ForTeam,
			&data.Sport,
			&data.LocLat,
			&data.LocLng,
			&data.Date)

		oppID, _ := strconv.ParseInt(data.Opposition, 10, 64)
		forID, _ := strconv.ParseInt(data.ForTeam, 10, 64)
		data.Opposition = GetTeamNameFromTeamID(oppID)
		data.ForTeam = GetTeamNameFromTeamID(forID)

		data.IsHome = false

		teamAway = append(teamAway, data)
	}

	// Initialise the json response for the end result
	var teamFixtures []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamFixtures)
	merge(&teamHome, &teamAway, &teamFixtures)

	j,_ := json.Marshal(teamFixtures)        // Convert the list of DB hits to a JSON
	// fmt.Println("Upcoming>>>")           // Write the result to the console
	// fmt.Println(string(j))           // Write the result to the console
	// fmt.Println("<<<Upcoming")           // Write the result to the console
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})


// Get the fixtures for the user specified in the url
var GetUserFixtures = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username=name)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := strings.Split(getquery, "=")[1]

	// Obtain userID
	userID := GetUserIDFromUsername(username)

	// Build queries
	ordering := "ORDER BY date ASC"
	commonQueryFields := "sport, loc_lat, loc_lng, date, score_home, score_away"
	tables := "previous_fixtures JOIN team_members"

	homeFields := fmt.Sprintf("away_id, home_id, %s", commonQueryFields)
	homeTableJoinCond := "home_id=team_id"

	awayFields := fmt.Sprintf("home_id, away_id, %s", commonQueryFields)
	awayTableJoinCond := "away_id=team_id"

	// Common text used in initialising JSON responses
	var jsonText = []byte(`[]`)

	homeQuery := fmt.Sprintf("SELECT %s FROM %s ON %s WHERE user_id=%d %s",
														homeFields, tables, homeTableJoinCond, userID, ordering)
	awayQuery := fmt.Sprintf("SELECT %s FROM %s ON %s WHERE user_id=%d %s",
														awayFields, tables, awayTableJoinCond, userID, ordering)

	// Run the query for home games
	rows, err := Database.Query(homeQuery)
	CheckErr(err)

	// Initialise the json response for all home games
	var teamHome []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamHome)

	// Add every database hit to the result
	for rows.Next() {
		data := Fixture{}
		err = rows.Scan(
			&data.Opposition,
			&data.ForTeam,
			&data.Sport,
			&data.LocLat,
			&data.LocLng,
			&data.Date,
			&data.ScoreHome,
			&data.ScoreAway)

		oppID, _ := strconv.ParseInt(data.Opposition, 10, 64)
		forID, _ := strconv.ParseInt(data.ForTeam, 10, 64)
		data.Opposition = GetTeamNameFromTeamID(oppID)
		data.ForTeam = GetTeamNameFromTeamID(forID)

		data.IsHome = true

		teamHome = append(teamHome, data)
	}

	rows, err = Database.Query(awayQuery)
	CheckErr(err)

	// Initialise the json response for all away games
	var teamAway []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamAway)

	// Add every database hit to the result
	for rows.Next() {
		data := Fixture{}
		err = rows.Scan(
			&data.Opposition,
			&data.ForTeam,
			&data.Sport,
			&data.LocLat,
			&data.LocLng,
			&data.Date,
			&data.ScoreHome,
			&data.ScoreAway)

		oppID, _ := strconv.ParseInt(data.Opposition, 10, 64)
		forID, _ := strconv.ParseInt(data.ForTeam, 10, 64)
		data.Opposition = GetTeamNameFromTeamID(oppID)
		data.ForTeam = GetTeamNameFromTeamID(forID)

		data.IsHome = false

		teamAway = append(teamAway, data)
	}

	// Initialise the json response for the end result
	var teamFixtures []Fixture
	err = json.Unmarshal([]byte(jsonText), &teamFixtures)
	merge(&teamHome, &teamAway, &teamFixtures)

	j,_ := json.Marshal(teamFixtures)        // Convert the list of DB hits to a JSON
	// fmt.Println("Prev>>>")           // Write the result to the console
	// fmt.Println(string(j))           // Write the result to the console
	// fmt.Println("<<<Prev")           // Write the result to the console
	fmt.Fprintln(writer, string(j)) // Write the result to the sender})
})


// Merge two lists of fixtures
func merge(list1 *[]Fixture, list2 *[]Fixture, result *[]Fixture) {
	var i, k int // Track positions in arrays

	for (i < len(*list1) || k < len(*list2)) {
		if (i >= len(*list1)) {
			*result = append(*result, (*list2)[k])
			k++
		} else if (k >= len(*list2)) {
			*result = append(*result, (*list1)[i])
			i++
		} else {
			if (strings.Compare((*list1)[i].Date, (*list2)[k].Date) >= 0) {
				*result = append(*result, (*list1)[i])
				i++
			} else {
				*result = append(*result, (*list2)[k])
				k++
			}
		}
	}
}


// Get the availability for the user specified in the url
var GetUserAvailability = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username=name)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := (strings.Split(getquery, "=")[1])

	// Obtain userID
	userID := GetUserIDFromUsername(username)

  // Run query
	daysFields := "mon, tues, weds, thurs, fri, sat, sun"
	query := fmt.Sprintf("SELECT %s FROM user_avail WHERE user_id=%d;",
		                    daysFields, userID)
	rows, err := Database.Query(query)
	CheckErr(err)

	// Initialise the json response for the result
	var result [7]int
	var jsonText = []byte(`[]`)
	err = json.Unmarshal([]byte(jsonText), &result)

	// Add the *only* database hit to the result
	if (rows.Next()) {
		err = rows.Scan(
			&result[0],
			&result[1],
			&result[2],
			&result[3],
			&result[4],
			&result[5],
			&result[6])
	} else {
		fmt.Println("Error no avail DB hit for ", username)
	}

	j,_ := json.Marshal(result) // Convert the list of DB hits to a JSON
	// fmt.Println(string(j))
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})


// Update the user availability for the user and values specified in the url
var UpdateUserAvailability = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain the bitmaps (query is of the form ?username=name&fst=x&snd=y)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := strings.Split((strings.Split(getquery, "=")[1]), "&")[0]

	monString := strings.Split((strings.Split(getquery, "=")[2]), "&")[0]
	monBitmap, _ := strconv.ParseInt(monString, 10, 64)

	tuesString := strings.Split((strings.Split(getquery, "=")[3]), "&")[0]
	tuesBitmap, _ := strconv.ParseInt(tuesString, 10, 64)

	wedsString := strings.Split((strings.Split(getquery, "=")[4]), "&")[0]
	wedsBitmap, _ := strconv.ParseInt(wedsString, 10, 64)

	thursString := strings.Split((strings.Split(getquery, "=")[5]), "&")[0]
	thursBitmap, _ := strconv.ParseInt(thursString, 10, 64)

	friString := strings.Split((strings.Split(getquery, "=")[6]), "&")[0]
	friBitmap, _ := strconv.ParseInt(friString, 10, 64)

	satString := strings.Split((strings.Split(getquery, "=")[7]), "&")[0]
	satBitmap, _ := strconv.ParseInt(satString, 10, 64)

	sunString := (strings.Split(getquery, "=")[8])
	sunBitmap, _ := strconv.ParseInt(sunString, 10, 64)

	// Obtain userID
	userID := GetUserIDFromUsername(username)

	// Run query
	fields := fmt.Sprintf("mon=%d, tues=%d, weds=%d, thurs=%d, fri=%d, sat=%d, sun=%d",
												 monBitmap, tuesBitmap, wedsBitmap, thursBitmap, friBitmap, satBitmap, sunBitmap)
	query := fmt.Sprintf("UPDATE user_avail SET %s WHERE user_id=%d",
												fields, userID)

	_, err = Database.Query(query)
	CheckErr(err)

	if err == nil {
		fmt.Fprintln(writer, "success") // Write the result to the sender
	} else {
		fmt.Fprintln(writer, "fail") // Write the result to the sender
	}

	recalculateUsersTeamAvailabilities(userID)
})


func recalculateUsersTeamAvailabilities(userID int) {
	query := fmt.Sprintf("SELECT team_id FROM team_members WHERE user_id=%d;", userID)

	rows, err := Database.Query(query)
	CheckErr(err)

	for (rows.Next()) {
		var team_id int
		rows.Scan(&team_id)
		RecalculateTeamAvailability(team_id)
	}
}


func RecalculateTeamAvailability(team_id int) {
	fields := "mon, tues, weds, thurs, fri, sat, sun"
		query := fmt.Sprintf("SELECT %s FROM user_avail NATURAL INNER JOIN team_members WHERE team_id=%d;", fields, team_id)

	rows, err := Database.Query(query)
	CheckErr(err)

	var totalMon   int64 = 0xFFFFFFFF
	var totalTues  int64 = 0xFFFFFFFF
	var totalWeds  int64 = 0xFFFFFFFF
	var totalThurs int64 = 0xFFFFFFFF
	var totalFri   int64 = 0xFFFFFFFF
	var totalSat   int64 = 0xFFFFFFFF
	var totalSun   int64 = 0xFFFFFFFF

	for (rows.Next()) {
		var holderMon   int64
		var holderTues  int64
		var holderWeds  int64
		var holderThurs int64
		var holderFri   int64
		var holderSat   int64
		var holderSun   int64

		rows.Scan(
			&holderMon,
			&holderTues,
			&holderWeds,
			&holderThurs,
			&holderFri,
			&holderSat,
			&holderSun)

		totalMon   = totalMon   & holderMon
		totalTues  = totalTues  & holderTues
		totalWeds  = totalWeds  & holderWeds
		totalThurs = totalThurs & holderThurs
		totalFri   = totalFri   & holderFri
		totalSat   = totalSat   & holderSat
		totalSun   = totalSun   & holderSun
	}

	fields = fmt.Sprintf("mon=%d, tues=%d, weds=%d, thurs=%d, fri=%d, sat=%d, sun=%d",
		                     totalMon, totalTues, totalWeds, totalThurs, totalFri, totalSat, totalSun)
	query = fmt.Sprintf("UPDATE team_avail SET %s WHERE team_id=%d",
											  fields, team_id)

	_, err = Database.Query(query)
	CheckErr(err)
}


// Update the user location for the user and values specified in the url
var UpdateUserLocation = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain the new location (query is of the form ?username=name&lat=x&lng=y)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	query := strings.Split(getquery, "&")

	fields := make([]string, len(query))
	for index, element:= range query{
		fields[index] = strings.Split(element, "=")[1]
	}

	username := fields[0]
	lat,err := strconv.ParseFloat(fields[1], 64)
	CheckErr(err)
	lng,err := strconv.ParseFloat(fields[2], 64)
	CheckErr(err)

	var userID int = GetUserIDFromUsername(username)

	// Run query
	dbfields := fmt.Sprintf("loc_lat=%f, loc_lng=%f", lat, lng)
	dbquery := fmt.Sprintf("UPDATE users SET %s WHERE user_id=%d",
											  dbfields, userID)

	_, err = Database.Query(dbquery)
	CheckErr(err)

	if err == nil {
		fmt.Fprintln(writer, "success") // Write the result to the sender
	} else {
		fmt.Fprintln(writer, "fail") // Write the result to the sender
	}

	recalculateUsersTeamLocations(userID);
})


func recalculateUsersTeamLocations(userID int) {
	query := fmt.Sprintf("SELECT team_id FROM team_members WHERE user_id=%d;", userID)

	rows, err := Database.Query(query)
	CheckErr(err)

	for (rows.Next()) {
		var team_id int
		rows.Scan(&team_id)
		RecalculateTeamLocation(team_id)
	}
}


func RecalculateTeamLocation(team_id int) {
	query := fmt.Sprintf("SELECT loc_lat, loc_lng FROM users NATURAL INNER JOIN team_members WHERE team_id=%d;", team_id)

	rows, err := Database.Query(query)
	CheckErr(err)

	var totalLat float64 = 0.0
	var totalLong float64 = 0.0
	var cnt float64 = 0

	for (rows.Next()) {
		var latHolder float64
		var longHolder float64
		rows.Scan(&latHolder, &longHolder)

		totalLat = totalLat + latHolder
		totalLong = totalLong + longHolder
		cnt = cnt + 1.0
	}

	query = fmt.Sprintf("UPDATE team_locations SET loc_lat=%f, loc_lng=%f WHERE team_id=%d;",
	                    (totalLat / cnt), (totalLong / cnt), team_id)

	_, err = Database.Query(query)
	CheckErr(err)
}
