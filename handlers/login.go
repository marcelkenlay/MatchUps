package handlers

import (
    "fmt"
    "log"
    "golang.org/x/crypto/bcrypt"

    "net/url"
    "strings"
    "net/http"
    _ "github.com/lib/pq"
    "encoding/json"
    . "../utils"
    "strconv"
)

type UserInfoInput struct {
	Username  string
	Name      string
	Dob       string
  LocLat    float64
	LocLng    float64
  Pwd       string
}

type UserLoginAttempt struct {
  Username  string
  Password  string
}

type UserLoginReturn struct {
  Error   string
  UserID  int
}

var AddUserInfo = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
  var userInfo UserInfoInput
  err := decoder.Decode(&userInfo)

  if err != nil {
      panic(err)
			defer request.Body.Close()
  }

  var hashed_pwd = HashPassword([]byte(userInfo.Pwd))

  // Run query to add user to DB
  fields := "username, name, dob, loc_lat, loc_lng, pwd_hash, score"
  query := fmt.Sprintf("INSERT INTO users (%s) VALUES('%s', '%s', '%s', %f, %f, '%s', '%s');",
							fields, userInfo.Username, userInfo.Name, userInfo.Dob,
              userInfo.LocLat, userInfo.LocLng, hashed_pwd, "100")
  fmt.Print(query)
  _, err = Database.Query(query)
  CheckErr(err)

  var userID = GetUserIDFromUsername(userInfo.Username)

  // Run query to add user's default availability to DB
  query = fmt.Sprintf("INSERT INTO user_avail VALUES (%d);", userID)
  _, err = Database.Query(query)
  CheckErr(err)

  fmt.Fprintln(writer, userID)
})



var GetLoginSuccess = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
  decoder := json.NewDecoder(request.Body)
  var userLoginAttempt UserLoginAttempt
  err := decoder.Decode(&userLoginAttempt)
  if err != nil {
      panic(err)
      defer request.Body.Close()
  }
	// Run query
  query := fmt.Sprintf("SELECT user_id, pwd_hash FROM users WHERE username='%s';", userLoginAttempt.Username)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	var userLoginReturn UserLoginReturn
  var pwd_hash string
	err = rows.Scan(
    &userLoginReturn.UserID,
		&pwd_hash)
  CheckErr(err)

  if (err != nil) {
    // If error then no entry was found in the database for the username given
    userLoginReturn.UserID = -1
    userLoginReturn.Error = "User Not Found"
  } else if (!ComparePasswords(pwd_hash, []byte(userLoginAttempt.Password))) {
    // If compare passwords returns false then we have an incorrect password attempt
    userLoginReturn.UserID = -1
    userLoginReturn.Error = "Incorrect Password"
  } else {
    // ComparePasswords returned true, username and password therefore valid
    userLoginReturn.Error = "none"
  }

  result, err := json.Marshal(userLoginReturn)
  CheckErr(err)
  fmt.Fprintln(writer, string(result))
})

var DoesMatchingUserExist = http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
	// Obtain username (query is of the form ?username)
	getquery, err := url.QueryUnescape(request.URL.RawQuery)
	username := strings.Split(getquery, "=")[1]

	// Run query

  query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE UPPER(username)='%s';", strings.ToUpper(username))
  // fmt.Println(query)
  rows, err := Database.Query(query)
  CheckErr(err)

	// Add the only database hit to the result
	rows.Next()
	var count string
	err = rows.Scan(
		&count)

  num, err := strconv.Atoi(count)
  match_exists := num > 0
  fmt.Fprintln(writer, match_exists)
})


func HashPassword(password []byte) string {
	bytes, err := bcrypt.GenerateFromPassword(password, 12)
  CheckErr(err)
	return string(bytes)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
  byteHash := []byte(hashedPwd)
  err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
  if (err != nil) {
      log.Println(err)
      return false
  }
  return true
}
