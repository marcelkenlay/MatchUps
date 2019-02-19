package handlers

import (
	. "./utils"
	. "./users"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"net/url"
	"strings"
)

var AddUserInfo = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var userInfo UserInfoInput
	err := decoder.Decode(&userInfo)

	if err != nil {
		panic(err)
	}

	userSessionCookie := InsertUserIntoDB(userInfo, userInfo.Pwd)

	InsertUserDefaultAvailIntoDB(userSessionCookie)

	result, err := json.Marshal(userSessionCookie)
	CheckErr(err)
	_, _ = fmt.Fprintln(writer, string(result))
})

var GetLoginSuccess = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var userLoginAttempt UserLoginAttempt
	err := decoder.Decode(&userLoginAttempt)
	if err != nil {
		panic(err)
	}

	userLoginReturn := CheckUserLogin(userLoginAttempt)

	result, err := json.Marshal(userLoginReturn)
	CheckErr(err)
	_, _ = fmt.Fprintln(writer, string(result))
})

var DoesMatchingUserExist = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, _ := url.QueryUnescape(request.URL.RawQuery)
	username := strings.Split(getquery, "=")[1]

	if username == "" {
		_, _ = fmt.Fprintln(writer, true)
	}
	matchExists := IsUsernameInUse(username)

	_, _ = fmt.Fprintln(writer, matchExists)
})