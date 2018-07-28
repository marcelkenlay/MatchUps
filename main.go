package main

import (
	"net/http"
  "log"
  _ "github.com/lib/pq"
	"github.com/gorilla/mux"
	. "./handlers"
)

func main() {
	ConnectToDatabase()
	defer Database.Close()

	SetUpPusher()

	r := mux.NewRouter()

	// Profile page
	r.Handle("/getuserinfo", GetUserInfo).Methods("GET")
	r.Handle("/getuserfixtures", GetUserFixtures).Methods("GET")
	r.Handle("/getuseravail", GetUserAvailability).Methods("GET")
	r.Handle("/updateuseravail", UpdateUserAvailability).Methods("GET")
	r.Handle("/getuserupcoming", GetUserUpcoming).Methods("GET")
	r.Handle("/updateuserloc", UpdateUserLocation).Methods("GET")

	// Chat Page
	r.Handle("/getChats",GetChats).Methods("GET")
	r.Handle("/getChatMessages", GetChatMessages).Methods("GET")
	r.Handle("/addMessage", AddMessage).Methods("POST")
	r.Handle("/getTeamMembers", GetTeamMembers).Methods("GET")
	r.Handle("/getTeamNames", GetTeamNames).Methods("GET")

	// Chat Page -- Promoted fixtures
	r.Handle("/getPromotedFixtures", GetPromotedFixtures).Methods("GET")
	r.Handle("/getUpvoteTally", GetUpvoteTally).Methods("GET")
	r.Handle("/getDownvoteTally", GetDownvoteTally).Methods("GET")
	r.Handle("/getVoteStatus", GetVoteStatus).Methods("GET")
	r.Handle("/addUpvote", AddUpvote).Methods("POST")
	r.Handle("/addDownvote", AddDownvote).Methods("POST")
	r.Handle("/removeUpvote", RemoveUpvote).Methods("POST")
	r.Handle("/removeDownvote", RemoveDownvote).Methods("POST")
	r.Handle("/acceptAdvert", AcceptAdvert).Methods("POST")
	r.Handle("/declineAdvert", DeclineAdvert).Methods("POST")

	// Chat Page -- Right Panel
	r.Handle("/getupcominggame", GetFixtureDetails).Methods("GET")
	r.Handle("/getsubmittedscore", GetSubmittedScore).Methods("GET")
	r.Handle("/acceptsubmittedscore", AcceptSubmittedScore).Methods("GET")
	r.Handle("/rejectsubmittedscore", RejectSubmittedScore).Methods("GET")
	r.Handle("/submitscore", SubmitScore).Methods("GET")
	r.Handle("/getprevgame", GetPreviousGame).Methods("GET")

	// Matchmaking
	r.HandleFunc("/matchmaking", GetMatchmaking).Methods("GET")

	// Team Page
	r.Handle("/getTeams/{user_id}", GetTeams).Methods("GET")
	r.Handle("/getInvitations/{username}", GetInvitations).Methods("GET")
	r.Handle("/addUserToTeam/{username}/{teamname}",AddPlayerToTeam).Methods("POST")
	r.Handle("/deleteInvitation/{username}/{teamname}",DeleteInvitation).Methods("DELETE")

	//Create Team Page
	r.Handle("/getUsernameMatches", GetUsernameMatches).Methods("GET")
	r.Handle("/createTeam", AddTeam).Methods("POST")
	r.Handle("/sendInvitations", SendInvitations).Methods("POST")

	// Login
	r.Handle("/addUserInfo", AddUserInfo).Methods("POST")
	r.Handle("/checkLogin", GetLoginSuccess).Methods("POST")
	r.Handle("/doesMatchingUserExist", DoesMatchingUserExist).Methods("GET")

	// Root
	http.Handle("/", r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./build/")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
