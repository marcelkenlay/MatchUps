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
const USERS_TABLE = "users"
const USER_ID_COL = "user_id"
const USERNAME_COL = "username"
const NAME_COL = "name"
const DOB_COL = "dob"
const LOC_LAT_COL = "loc_lat"
const LOC_LNG_COL = "loc_lng"
const PWD_HASH_COL = "pwd_hash"
const SCORE_COL = "score"

func InsertUserIntoDB(userInfo UserInfoInput, password string) UserSessionCookie {
	hashedPwd := HashPassword(password)

	columns := []string{USERNAME_COL, NAME_COL, DOB_COL, LOC_LAT_COL, LOC_LNG_COL, PWD_HASH_COL}
	args := []interface{}{
		userInfo.Username, userInfo.Name, userInfo.Dob,
		userInfo.LocLat, userInfo.LocLng, hashedPwd}

	id, _ := InsertRowIntoTableAndRetreiveVal(USERS_TABLE, columns, args, USER_ID_COL)
	idNum, _ := strconv.Atoi(id)
	return GenerateSessionCookie(idNum)
}

func CheckUserLogin(userLoginAttempt UserLoginAttempt) UserLoginReturn {

	columns := []string{USER_ID_COL, PWD_HASH_COL}
	conditions := []Condition{SingleValCondition(fmt.Sprintf("%s = ?", USERNAME_COL), userLoginAttempt.Username)}

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
	conditions := []Condition{SingleValColEqCondition(USERNAME_COL, username)}

	count, err := CountMatchingRows(USERS_TABLE, conditions)

	if err != nil {
		log.Println("Failed to query for count result")
		return true
	}

	return count > 0
}

func GetUserLocationInfo(userSession UserSessionCookie) string {
	userId := GetUserIdFromSession(userSession)

	locLat, locLng := GetUserLatLngFromId(userId)

	locLatS, locLngS := fmt.Sprintf("%f", locLat) , fmt.Sprintf("%f", locLng)

	return GetAddressForLatLng(locLatS, locLngS)
}

func GetUserLatLngFromId(userId int) (float64, float64) {
	columns := []string{LOC_LAT_COL, LOC_LNG_COL}

	conditions := []Condition{SingleValColEqCondition(USER_ID_COL, userId)}

	row, err := SelectRowFromTable(USERS_TABLE, columns, conditions)
	CheckErr(err)

	var locLng, locLat float64

	err = row.Scan(&locLat, &locLng)

	return locLat, locLng
}

func GetIdsForUsernames(usernames []string) (ids []int) {
	columns := []string{USER_ID_COL}

	var args []interface{}
	for _, username := range usernames {
		args = append(args, username)
	}

	conditions := []Condition{ColInSetCondition(USERNAME_COL, args)}

	rows, _ := SelectRowsFromTable(USERS_TABLE, columns, conditions)

	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}
