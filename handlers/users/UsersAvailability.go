package users


import (
	. "../db"
)

var USER_AVAIL_TABLE = "user_availability"
var USER_ID = "user_id"

func InsertUserDefaultAvailIntoDB(userSession UserSessionCookie) {

	userId := GetUserIdFromSession(userSession)

	columns := []string{USER_ID}
	vals :=   []interface{}{userId}

	_ = InsertRowIntoTable(USER_AVAIL_TABLE, columns, vals)
}
