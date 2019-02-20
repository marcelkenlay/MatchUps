package users


import (
	. "../db"
)

var USER_AVAIL_TABLE = "user_availability"

func InsertUserDefaultAvailIntoDB(userSession UserSessionCookie) {

	userId := GetUserIdFromSession(userSession)

	columns := []string{USER_ID_COL}
	vals :=   []interface{}{userId}

	_ = InsertRowIntoTable(USER_AVAIL_TABLE, columns, vals)
}
