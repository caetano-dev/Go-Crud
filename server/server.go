package server

import (
	"database/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser creates users in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.Write([]byte("Error reading request body"))
		return
	}

	var user user
	if error = json.Unmarshal(requestBody, &user); error != nil {
		w.Write([]byte("Error converting user"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		w.Write([]byte("Error connecting to DB"))
		return
	}

	statement, error := db.Prepare("insert into users (name, email) values (?, ?)")
	if error != nil {
		w.Write([]byte("Error loading statement"))
		return
	}

	defer statement.Close()

	insertion, error := statement.Exec(user.Name, user.Email)

	if error != nil {
		w.Write([]byte("Insertion Error"))
		return
	}

	insertedID, error := insertion.LastInsertId()
	if error != nil {

		w.Write([]byte("Error obtaining id"))
		return

	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Success inserting user! ID: %d", insertedID)))
}

//FetchUsers fetches users in the database
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	db, error := database.Connect()

	if error != nil {

		w.Write([]byte("Error connecting to database"))
		return
	}

	defer db.Close()
	lines, error := db.Query("select * from users")

	if error != nil {

		w.Write([]byte("Error fetching users"))
		return
	}

	defer lines.Close()
	var users []user
	for lines.Next() {
		var user user
		if error := lines.Scan(&user.ID, &user.Name, &user.Email); error != nil {
			w.Write([]byte("Error scanning"))
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(w).Encode(users); error != nil {

		w.Write([]byte("Error converting to json"))
		return
	}
}

//FetchUser fetches an user
func FetchUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, error := strconv.ParseInt(params["id"], 10, 32)
	if error != nil {
		w.Write([]byte("Error converting parameter to int"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		w.Write([]byte("Error connecting to database"))
		return
	}

	line, error := db.Query("select * from users where id = ?", ID)

	if error != nil {
		w.Write([]byte("Error fetching user"))
		return
	}

	var user user
	if line.Next() {
		if error := line.Scan(&user.ID, &user.Name, &user.Email); error != nil {

			w.Write([]byte("Error scanning user"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(w).Encode(user); error != nil {
		w.Write([]byte("Error converting user to json"))
		return
	}

}

//UpdateUser updates an user in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, error := strconv.ParseInt(params["id"], 10, 32)
	if error != nil {
		w.Write([]byte("Error converting parameter to int"))
		return
	}

	requestBody, error := ioutil.ReadAll(r.Body)

	if error != nil {
		w.Write([]byte("Error reading request body"))
		return
	}

	var user user

	if error = json.Unmarshal(requestBody, &user); error != nil {

		w.Write([]byte("Error converting user to struct"))
		return
	}

	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Error connecting to DB"))
		return
	}

	defer db.Close()

	statement, error := db.Prepare("update users set name = ?, email = ? where id = ?")

	if error != nil {
		w.Write([]byte("Error creating statement"))
		return
	}

	defer statement.Close()

	if _, error := statement.Exec(user.Name, user.Email, ID); error != nil {
		w.Write([]byte("Error updating user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
