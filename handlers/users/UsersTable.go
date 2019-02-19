package users

import (
	. "../db"
	. "../request_utils"
	. "../utils"
	"fmt"
	"log"
	"strconv"
)

type UserLoginAttempt struct {
	Username string
	Password string
}

type UserLoginReturn struct {
	Error             string
	UserSessionCookie UserSessionCookie
}

type UserInfoInput struct {
	Username string
	Name     string
	Dob      string
	LocLat   float64
	LocLng   float64
	Pwd      string
}

// Users table
var USERS_TABLE = "users"
var ID = "user_id"
var USERNAME = "username"
var NAME = "name"
var DOB = "dob"
var LOC_LAT = "loc_lat"
var LOC_LNG = "loc_lng"
var PWD_HASH = "pwd_hash"
var SCORE = "score"

func InsertUserIntoDB(userInfo UserInfoInput, password string) UserSessionCookie {
	hashedPwd := HashPassword(password)

	columns := []string{USERNAME, NAME, DOB, LOC_LAT, LOC_LNG, PWD_HASH}
	args := []interface{}{
		userInfo.Username, userInfo.Name, userInfo.Dob,
		userInfo.LocLat, userInfo.LocLng, hashedPwd}

	id, _ := InsertRowIntoTableAndRetreiveVal(USERS_TABLE, columns, args, USER_ID)
	idNum, _ := strconv.Atoi(id)
	return GenerateSessionCookie(idNum)
}

func CheckUserLogin(userLoginAttempt UserLoginAttempt) UserLoginReturn {

	columns := []string{USER_ID, PWD_HASH}
	conditions := []Condition{SingleValCondition(fmt.Sprintf("%s = ?", USERNAME), userLoginAttempt.Username)}

	row, err := SelectRowFromTable(USERS_TABLE, columns, conditions)

	var userLoginReturn UserLoginReturn
	var correctPwdHash string
	var userId int

	err = row.Scan(&userId, &correctPwdHash)

	if err != nil {
		// If error then no entry was found in the database for the username given
		userLoginReturn.UserSessionCookie = UserSessionCookie{}
		log.Println(err)
		userLoginReturn.Error = "User Not Found"
	} else if !ComparePasswords(correctPwdHash, []byte(userLoginAttempt.Password)) {
		// If compare passwords returns false then we have an incorrect password attempt
		userLoginReturn.UserSessionCookie = UserSessionCookie{}
		userLoginReturn.Error = "Incorrect Password"
	} else {
		// ComparePasswords returned true, username and password therefore valid
		userLoginReturn.UserSessionCookie = GenerateSessionCookie(userId)
		userLoginReturn.Error = "none"
	}
	return userLoginReturn
}

func IsUsernameInUse(username string) bool {
	conditions := []Condition{SingleValColEqCondition(USERNAME, username)}
	columns := []string{"COUNT(*)"}

	row, err := SelectRowFromTable(USERS_TABLE, columns, conditions)

	if err != nil {
		log.Println("Failed to query for count result")
		return true
	}

	var count string
	err = row.Scan(&count)

	if err != nil {
		log.Println("Failed to extract count result")
		return true
	}

	num, _ := strconv.Atoi(count)

	return num > 0
}

func GetUserLocationInfo(userSession UserSessionCookie) string {
	userId := GetUserIdFromSession(userSession)

	columns := []string{LOC_LAT, LOC_LNG}

	conditions := []Condition{SingleValColEqCondition(USER_ID, userId)}

	row, err := SelectRowFromTable(USERS_TABLE, columns, conditions)
	CheckErr(err)

	var locLng, locLat string

	err = row.Scan(&locLat, &locLng)

	return GetAddressForLatLng(locLat, locLng)
}
