package users

import (
	. "../db"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	mr "math/rand"
	"net/http"
	"strconv"
	"time"
)

type UserSessionCookie struct {
	ID   int
	Hash string
}

// User Sessions table
var USER_SESSIONS_TABLE = "user_sessions"
var COOKIE_HASH = "cookie_hash"
var SESSION_ID = "session_id"

func BuildUserSessionFromRequest(request *http.Request) UserSessionCookie {
	var sessionId, sessionHash string
	for _, cookie := range request.Cookies() {
		switch cookie.Name {
		case "UserSessionId":
			sessionId = cookie.Value
			break
		case "UserSessionHash":
			sessionHash = cookie.Value
			break
		}
	}
	if sessionId == "" || sessionHash == "" {
		log.Println("Error reading session from request")
	}
	sessionIdI, _ := strconv.Atoi(sessionId)
	return UserSessionCookie{ID: sessionIdI, Hash: sessionHash}
}

func GetUserIdFromSession(userSession UserSessionCookie) int {
	columns := []string{USER_ID_COL, COOKIE_HASH}

	conditions := []Condition{SingleValCondition(fmt.Sprintf("%s = ?", SESSION_ID), userSession.ID)}

	row, _ := SelectRowFromTable(USER_SESSIONS_TABLE, columns, conditions)

	var userId int
	var correctCookieHash string

	_ = row.Scan(&userId, &correctCookieHash)

	if ComparePasswords(correctCookieHash, []byte(userSession.Hash)) {
		return userId
	}
	return -1
}

func InsertUserSession(userId int, sessionCookieHash string) (int, error) {
	columns := []string{USER_ID_COL, COOKIE_HASH}
	args := []interface{}{userId, sessionCookieHash}

	sessionId, err := InsertRowIntoTableAndRetreiveVal(USER_SESSIONS_TABLE, columns, args, SESSION_ID)

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(sessionId)

}

func GenerateSessionCookie(userId int) UserSessionCookie {
	count := 1
	for count < 10 {
		b := String(16)
		hash := HashPassword(b)
		sessionId, err := InsertUserSession(userId, hash)
		if err == nil {
			log.Println(string(b))
			return UserSessionCookie{sessionId, string(b)}
		}
		count++
	}
	panic("cannot create a user session")
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = mr.New(mr.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func HashByteArray(ary []byte) string {
	hash, err := bcrypt.GenerateFromPassword(ary, 12)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(hash)
}

func HashPassword(password string) string {
	bytePwd := []byte(password)
	return HashByteArray(bytePwd)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
