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

type MessageInfo struct {
	Chat 			 string
	SenderID	 int
	SenderName string
	Message		 string
	TimeSent   string
}

type ChannelMessage struct {
	SenderID	 int
	SenderName string
	Message		 string
	TimeSent   string
}

type ChatInfo struct {
	FixtureID			int
	UserTeamID		int
	OppID					int
	UserTeamName	string
	OppName				string
}

var GetChats = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain user id (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	user_id := (strings.Split(getquery, "=")[1])

	// Define struct to hold info on all chats to be returned
	var chats []ChatInfo
	// Define struct to hold info on all teams user belongs to
	var team_ids []int

	// Obtain team information
	columns := "team_names.team_id, team_names.team_name"
	join := "team_members INNER JOIN team_names ON team_names.team_id=team_members.team_id"
  query := fmt.Sprintf("SELECT %s FROM %s where user_id=%s;", columns, join, user_id)
  rows, err := Database.Query(query)
  CheckErr(err)
	for rows.Next() {
		data := ChatInfo{}
		data.FixtureID = -1
		err = rows.Scan(
			&data.UserTeamID,
			&data.UserTeamName)
		chats = append(chats, data)
		team_ids = append(team_ids, data.UserTeamID)
	}

	// Obtain fixture information
	for _, team_id := range team_ids {
		// Obtain fixtures for team
		columns := "f.fixture_id, f.home_id, f.away_id, h.team_name, a.team_name"
		fix_table := "upcoming_fixtures f"
		first_join := "INNER JOIN team_names h ON f.home_id=h.team_id"
		second_join := "INNER JOIN team_names a ON f.away_id=a.team_id"
  	query := fmt.Sprintf("SELECT %s FROM %s %s %s WHERE home_id=%d OR away_id=%d;",
			 									columns, fix_table, first_join, second_join, team_id, team_id)
  	rows, err := Database.Query(query)
  	CheckErr(err)
		var home_id int
		var away_id int
		var home_name string
		var away_name string
		for rows.Next() {
			data := ChatInfo{}
			err = rows.Scan(
				&data.FixtureID,
				&home_id,
				&away_id,
				&home_name,
				&away_name)
			if (home_id == team_id) {
				data.UserTeamID = home_id;
				data.OppID = away_id;
				data.UserTeamName = home_name;
				data.OppName = away_name;
			} else {
				data.UserTeamID = away_id;
				data.OppID = home_id;
				data.UserTeamName = away_name;
				data.OppName = home_name;
			}
			chats = append(chats, data)
		}

	}

	j,_ := json.Marshal(chats) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

var GetChatMessages = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	chat_id := (strings.Split(getquery, "=")[1])
	chat_db_name := chat_id + "_messages"

	// Run query
	columns := "m.sender_id, u.name, m.message, m.time_sent"
	join := "INNER JOIN users AS u ON m.sender_id=u.user_id"
	order:= "ORDER BY m.time_sent ASC"

  query := fmt.Sprintf("SELECT %s FROM %s AS m %s %s",
		 						columns, chat_db_name, join, order)
	// fmt.Println(query)
  rows, err := Database.Query(query)
  CheckErr(err)

	var result []ChannelMessage

	for rows.Next() {
		data := ChannelMessage{}
		err = rows.Scan(
			&data.SenderID,
			&data.SenderName,
			&data.Message,
			&data.TimeSent)

		result = append(result, data)
	}

	j,_ := json.Marshal(result) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})

func escape(source string) string {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\'':
			flag = true
			escape = '\''
			break
		default:
		}
		if flag {
			desc[j] = '\''
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

var AddMessage = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
  var message MessageInfo
  err := decoder.Decode(&message)
  if err != nil {
      panic(err)
			defer request.Body.Close()
  }
	chat_name := message.Chat + "_messages"
	columns := "sender_id, message, time_sent"
  query := fmt.Sprintf("INSERT INTO %s (%s) VALUES('%d', '%s', LOCALTIMESTAMP);",
							chat_name, columns, message.SenderID, escape(message.Message))
  _, err = Database.Query(query)
  CheckErr(err)

	// Message to send to other subscribers to channel
	 channel_message := ChannelMessage{SenderID: message.SenderID,
			 SenderName: message.SenderName, Message: message.Message,
			 TimeSent: message.TimeSent}

	 // trigger an event on a channel, along with a data payload
	 PusherClient.Trigger(message.Chat, "message", channel_message)
})


var GetTeamMembers = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	team_id := strings.Split(getquery, "=")[1]

	// Run query
  query := fmt.Sprintf("SELECT team_members.user_id FROM team_members WHERE team_members.team_id=%s;", team_id)
  rows, err := Database.Query(query)
  CheckErr(err)

	var result []int
	// Add every database hit to the result
	for rows.Next() {
		var member int
		err = rows.Scan(&member)
		result = append(result, member)
	}

	j,_ := json.Marshal(result) // Convert the list of DB hits to a JSON
	fmt.Fprintln(writer, string(j)) // Write the result to the sender
})
