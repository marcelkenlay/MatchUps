package handlers

import (
	"net/http"
  _ "github.com/lib/pq"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	. "../utils"
)

type PromotedAdvertisement struct {
	AdID			 int
	TeamID     string
	Name       string
	StartTime  string
	EndTime  	 string
	LocLat	 	 float64
	LocLng	 	 float64
	Sport   	 string
	NumPlayers int
}

var GetUpvoteTally = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	params := strings.Split( strings.Split(getquery, "?")[0], "&" )
	team_id := strings.Split(params[0], "=")[1]
	advert_id := strings.Split(params[1], "=")[1]

	query := fmt.Sprintf("SELECT COUNT(*) FROM upvotes WHERE team_id=%s AND advert_id=%s",
		 						team_id, advert_id)
	// fmt.Println(query)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	var count string
	err = rows.Scan(&count)

	j,_ := json.Marshal(count) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

var GetDownvoteTally = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	params := strings.Split( strings.Split(getquery, "?")[0], "&" )
	team_id := strings.Split(params[0], "=")[1]
	advert_id := strings.Split(params[1], "=")[1]

	query := fmt.Sprintf("SELECT COUNT(*) FROM downvotes WHERE team_id=%s AND advert_id=%s",
		 						team_id, advert_id)
	// fmt.Println(query)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	var count int
	err = rows.Scan(&count)

	j,_ := json.Marshal(count) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

var GetVoteStatus = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	params := strings.Split( strings.Split(getquery, "?")[0], "&" )
	user_id := strings.Split(params[0], "=")[1]
	team_id := strings.Split(params[1], "=")[1]
	advert_id := strings.Split(params[2], "=")[1]

  // Check for instance of an upvote
	query := fmt.Sprintf("SELECT COUNT(*) FROM upvotes WHERE user_id=%s AND team_id=%s AND advert_id=%s",
		 						user_id, team_id, advert_id)
	// fmt.Println(query)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	var count int
	err = rows.Scan(&count)

  if (count > 0) {
    fmt.Fprintln(writer, "upvote") // Write the result to the sender
    return
  }

  // Check for instance of a downvote
  query = fmt.Sprintf("SELECT COUNT(*) FROM downvotes WHERE user_id=%s AND team_id=%s AND advert_id=%s",
		 						user_id, team_id, advert_id)
	// fmt.Println(query)
  rows, err = Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	err = rows.Scan(&count)

  if (count > 0) {
    fmt.Fprint(writer, "downvote") // Write the result to the sender
    return
  }

  fmt.Fprint(writer, "novote")

})

type Vote struct {
  UserID    int
  TeamID    int
	AdvertID	int
}

//todo set up MUX router to take url of user and team to add to database.
var AddUpvote = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var vote Vote
	err := decoder.Decode(&vote)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Insert into the upvote table
  query := fmt.Sprintf("INSERT INTO upvotes (user_id, team_id, advert_id) VALUES(%d, %d, %d);",
              vote.UserID, vote.TeamID, vote.AdvertID)
  _, err = Database.Query(query)
  CheckErr(err)

	// Push upvote addition
	channel := fmt.Sprintf("_team%d", vote.TeamID)
	message := fmt.Sprintf("ad_%d_add_up", vote.AdvertID)
	PusherClient.Trigger(channel, message, vote.UserID)

})


//todo set up MUX router to take url of user and team to add to database.
var AddDownvote = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var vote Vote
	err := decoder.Decode(&vote)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Insert into the upvote table
  query := fmt.Sprintf("INSERT INTO downvotes (user_id, team_id, advert_id) VALUES(%d, %d, %d);",
              vote.UserID, vote.TeamID, vote.AdvertID)
  _, err = Database.Query(query)
  CheckErr(err)

	channel := fmt.Sprintf("_team%d", vote.TeamID)
	message := fmt.Sprintf("ad_%d_add_dwn", vote.AdvertID)
	PusherClient.Trigger(channel, message, vote.UserID)

})


//todo set up MUX router to take url of user and team to add to database.
var RemoveUpvote = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var vote Vote
	err := decoder.Decode(&vote)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Delete from the upvote table
  query := fmt.Sprintf("DELETE FROM upvotes WHERE user_id=%d AND team_id=%d AND advert_id=%d;",
              vote.UserID, vote.TeamID, vote.AdvertID)
  _, err = Database.Query(query)
  CheckErr(err)

	channel := fmt.Sprintf("_team%d", vote.TeamID)
	message := fmt.Sprintf("ad_%d_rmv_up", vote.AdvertID)
	PusherClient.Trigger(channel, message, vote.UserID)
})


//todo set up MUX router to take url of user and team to add to database.
var RemoveDownvote = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var vote Vote
	err := decoder.Decode(&vote)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Delete from the downvote table
	query := fmt.Sprintf("DELETE FROM downvotes WHERE user_id=%d AND team_id=%d AND advert_id=%d;",
              vote.UserID, vote.TeamID, vote.AdvertID)
  _, err = Database.Query(query)
  CheckErr(err)

	channel := fmt.Sprintf("_team%d", vote.TeamID)
	message := fmt.Sprintf("ad_%d_rmv_dwn", vote.AdvertID)
	PusherClient.Trigger(channel, message, vote.UserID)

})

var GetPromotedFixtures = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	team_id := strings.Split(getquery, "=")[1]

	fields := "a.advert_id, a.team_id, tn.team_name, a.start_time, a.end_time, a.loc_lat, a.loc_lng, a.sport, a.num_players";
	advertisements := "advertisements AS a"
	first_join := "INNER JOIN promoted_fixtures AS pf ON a.advert_id=pf.advert_id"
	second_join := "INNER JOIN team_names AS tn ON a.team_id=tn.team_id"

	// Run query
  query := fmt.Sprintf("SELECT %s FROM %s %s %s WHERE pf.team_id=%s;",
												fields, advertisements, first_join, second_join, team_id)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Initialise the json response
	var jsonText = []byte(`[]`)
	var result []PromotedAdvertisement
	err = json.Unmarshal([]byte(jsonText), &result)

	// Add every database hit to the result
	for rows.Next() {
		data := PromotedAdvertisement{}
		err = rows.Scan(
			&data.AdID,
			&data.TeamID,
			&data.Name,
			&data.StartTime,
			&data.EndTime,
			&data.LocLat,
			&data.LocLng,
			&data.Sport,
      &data.NumPlayers)

		result = append(result, data)
	}

	j,_ := json.Marshal(result) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})


type Acceptance struct {
  AccepterID int
	AdID			 int
	HostID     int
	StartTime  string
	LocLat	 	 float64
	LocLng	 	 float64
	Sport   	 string
}

//todo set up MUX router to take url of user and team to add to database.
var AcceptAdvert = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var acceptance Acceptance
	err := decoder.Decode(&acceptance)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Delete from the promoted_fixtures table
	query := fmt.Sprintf("DELETE FROM promoted_fixtures WHERE advert_id=%d;", acceptance.AdID)
	_, err = Database.Query(query)
  CheckErr(err)

	// Delete from the advertisements table
	query = fmt.Sprintf("DELETE FROM advertisements WHERE advert_id=%d;", acceptance.AdID)
	_, err = Database.Query(query)
  CheckErr(err)

	// Add new fixture
	query = fmt.Sprintf("INSERT INTO upcoming_fixtures " +
											"(fixture_id, home_id, away_id, sport, loc_lat, loc_lng, date)" +
											" VALUES(%d, %d, %d, '%s', %f, %f, '%s');",
											acceptance.AdID, acceptance.HostID, acceptance.AccepterID, acceptance.Sport,
											acceptance.LocLat, acceptance.LocLng, acceptance.StartTime)
	fmt.Println(query)
	_, err = Database.Query(query)
	CheckErr(err)


	//Create message table for team
	table_name := fmt.Sprintf("_fixture%d_messages", acceptance.AdID)
	columns := "sender_id integer NOT NULL, message varchar(200) NOT NULL, Time_sent timestamp without time zone NOT NULL"
	query = fmt.Sprintf("CREATE TABLE %s (%s);", table_name, columns)
	fmt.Println(query)
	_, err = Database.Query(query)
	CheckErr(err)
})

type Declined struct {
  DeclinerID int
	AdID			 int
}

var DeclineAdvert = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var declined Declined
	err := decoder.Decode(&declined)
	if err != nil {
		panic(err)
		defer request.Body.Close()
	}

  // Delete from the promoted_fixtures table
	query := fmt.Sprintf("DELETE FROM promoted_fixtures WHERE advert_id=%d AND team_id=%d;",
		 declined.AdID, declined.DeclinerID)
	_, err = Database.Query(query)
  CheckErr(err)

})
