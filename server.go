package main

// TODO: Add authentication to all requests

import (
	"database/sql"
	"encoding/json"
	"go_chatservice/internal/lib"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-sql-driver/mysql"
)

// DB is the global database connection
var DB *sql.DB

// GetDB connects to the database and sets the DB variable
func GetDB() {
	// Create a MySQL configuration struct
	cfg := mysql.Config{
		User:   "chatservice",
		Passwd: "user",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "chatservice",
	}

	// Open a connection to the database
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Successfully connected to database!")
}

/**
 * Creates a new user and inserts it into the database.
 *
 * @param w http.ResponseWriter
 * @param body []byte
 * @return void
 */
func userNew(w http.ResponseWriter, body []byte) {
	// Unmarshal the body into a `User` struct.
	var user struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Check if the user data is non-default or valid.
	if user.Username == "" || user.Firstname == "" || user.Lastname == "" || user.Email == "" ||
		user.Password == "" {
		log.Debug("Error: User data is default")
		http.Error(w, "user data is default/empty", http.StatusBadRequest)
		return
	}

	// Check if the username is already taken.
	rows, err := DB.Query("SELECT * FROM User WHERE username = ?;", user.Username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	if rows.Next() {
		log.Debug("Error: Username already taken")
		http.Error(w, "username already taken", http.StatusBadRequest)
		return
	}

	// Insert the user into the database.
	stmt, err := DB.Prepare(
		"INSERT INTO User (username, first_name, last_name, email, password, joined_on) VALUES (?, ?, ?, ?, ?, ?);",
	)
	if err != nil {
		log.Error("Error preparing statement: ", err)
		http.Error(w, "can't prepare statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	stmt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Password, time.Now())

	log.Debug("inserted new user into database!")
}

/**
 * Gets user-info for supplied username
 *
 * @param w http.ResponseWriter
 * @param body []byte
 * @return void
 */
func userGet(w http.ResponseWriter, body []byte) {
	// Decode body into user
	var user lib.User
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Get username from request
	username := user.Username

	// Validate username
	if username == "" {
		log.Debug("Error: Username is empty")
		http.Error(w, "username is empty", http.StatusBadRequest)
		return
	}

	// Query database for user
	rows, err := DB.Query("SELECT * FROM User WHERE username = ?;", username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Check if user exists
	if !rows.Next() {
		log.Debug("Error: User not found")
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	// Scan user data from database
	var joined_raw string
	err = rows.Scan(
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&joined_raw,
	)
	if err != nil {
		log.Error("Error scanning row: ", err)
		http.Error(w, "can't scan row", http.StatusInternalServerError)
		return
	}

	// Parse joined date from database
	joined, err := time.Parse(time.DateTime, joined_raw)
	if err != nil {
		log.Error("Error parsing database date: ", err)
		http.Error(w, "can't parse database date", http.StatusInternalServerError)
		return
	}

	// Set joined date on user
	user.Joined = joined

	// Clear password from user
	user.Password = ""

	// Marshal user to JSON
	user_json, err := json.Marshal(user)
	if err != nil {
		log.Error("Error marshalling user: ", err)
		http.Error(w, "can't marshal user", http.StatusInternalServerError)
		return
	}

	// Get all chatgroups the user is in
	rows, err = DB.Query("SELECT chat_group_name FROM ChatGroupMemberShip WHERE user_name = ?;", username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var chat_groups []string
	// Iterate over all rows
	for rows.Next() {
		// Get chat_group_name from row
		var chat_group_name string
		err = rows.Scan(&chat_group_name)
		if err != nil {
			log.Error("Error scanning row: ", err)
			http.Error(w, "can't scan row", http.StatusInternalServerError)
			return
		}

		log.Debug("chat_group_name: ", chat_group_name)

		// Append chat_group_name
		chat_groups = append(chat_groups, chat_group_name)

	}

	// Marshal chat_groups to JSON
	chat_groups_json, err := json.Marshal(chat_groups)
	if err != nil {
		log.Error("Error marshalling chat_groups: ", err)
		http.Error(w, "can't marshal chat_groups", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Write(user_json)
	w.Write(chat_groups_json)
	w.Header().Set("Content-Type", "application/json")
	return
}

/**
 * Handles user requests (POST, GET).
 *
 * @param w The http.ResponseWriter.
 * @param req The http.Request.
 * @return void
 */
func userHandler(w http.ResponseWriter, req *http.Request) {
	// Read the request body.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Log the error and return an error response.
		log.Debug("Error reading body: ", err)
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}

	// Route the request based on the method.
	switch req.Method {
	case http.MethodPost:
		// Log the request and call the `userNew` function.
		log.Debug("/user POST")
		userNew(w, body)
	case http.MethodGet:
		// Log the request and call the `userGet` function.
		log.Debug("/user GET")
		userGet(w, body)
	}
}

/**
 * Creates a new chatgroup and inserts it into the database.
 *
 * @param w http.ResponseWriter
 * @param body []byte
 * @return void
 */
func chatgroupNew(w http.ResponseWriter, body []byte) {
	// TODO: Add authentication checking for the user who wants to create a chatgroup

	// Unmarshal body into Chatgroup struct
	var chatgroup lib.ChatGroup
	err := json.Unmarshal(body, &chatgroup)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Validate chatgroup data
	if chatgroup.Name == "" || chatgroup.CreatedBy == "" {
		log.Debug("Error: Chatgroup data is default")
		http.Error(w, "chatgroup data is default/empty", http.StatusBadRequest)
		return
	}

	// Check if chatgroup name is already taken
	rows, err := DB.Query("SELECT * FROM ChatGroup WHERE name = ?;", chatgroup.Name)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	if rows.Next() { // Chatgroup name already taken
		log.Debug("Error: Chatgroup name already taken")
		http.Error(w, "chatgroup name already taken", http.StatusBadRequest)
		return
	}

	// Check if user exists
	rows, err = DB.Query("SELECT * FROM User WHERE username = ?;", chatgroup.CreatedBy)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}
	if !rows.Next() { // User does not exist
		log.Debug("Error: User does not exist")
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	//TODO: check if user password is correct

	// Insert chatgroup into database
	stmt, err := DB.Prepare(
		"INSERT INTO ChatGroup (name, created_on, created_by) VALUES (?, ?, ?);",
	)
	if err != nil {
		log.Error("Error preparing statement: ", err)
		http.Error(w, "can't prepare statement", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(chatgroup.Name, time.Now(), chatgroup.CreatedBy)
	if err != nil {
		log.Error("Error executing statement: ", err)
		http.Error(w, "can't execute statement", http.StatusInternalServerError)
		return
	}

}

/**
 * Gets chatgroup-info for supplied name
 *
 * @param w http.ResponseWriter
 * @param body []byte
 * @return void
 */
func chatgroupGet(w http.ResponseWriter, body []byte) {
	// Decode body into placeholder chatgroup
	var placeholderChatGroup lib.ChatGroup
	err := json.Unmarshal(body, &placeholderChatGroup)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Get chatgroup name from placeholder chatgroup
	name := placeholderChatGroup.Name

	// Validate chatgroup name
	if name == "" {
		log.Debug("Error: Chatgroup name is empty")
		http.Error(w, "chatgroup name is empty", http.StatusBadRequest)
		return
	}

	// Query database for chatgroup
	rows, err := DB.Query("SELECT * FROM ChatGroup WHERE name = ?;", name)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Check if chatgroup was found
	if !rows.Next() {
		log.Debug("Error: Chatgroup not found")
		http.Error(w, "chatgroup not found", http.StatusBadRequest)
		return
	}

	// Initialize chatgroup struct
	var chatGroup lib.ChatGroup

	// Created_on as string (since we cant parse time.Time directly sadly)
	var createdRaw string

	// Scan row into chatgroup
	err = rows.Scan(
		&chatGroup.Name,
		&createdRaw,
		&chatGroup.CreatedBy,
	)

	if err != nil {
		log.Error("Error scanning row: ", err)
		http.Error(w, "can't scan row", http.StatusInternalServerError)
		return
	}

	// Convert created_on string to time.Time
	chatGroup.CreatedOn, err = time.Parse(time.DateTime, createdRaw)
	if err != nil {
		log.Error("Error parsing database date: ", err)
		http.Error(w, "can't parse database date", http.StatusInternalServerError)
		return
	}

	// Marshal chatgroup to JSON
	chatGroupJSON, err := json.Marshal(chatGroup)
	if err != nil {
		log.Error("Error marshalling chatgroup: ", err)
		http.Error(w, "can't marshal chatgroup", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Write(chatGroupJSON)
	w.Header().Set("Content-Type", "application/json")
}

/**
 * Handles chatgroup requests (POST, GET).
 *
 * @param w The http.ResponseWriter.
 * @param req The http.Request.
 * @return void
 */
func chatgroup_handler(w http.ResponseWriter, req *http.Request) {
	// Read request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// Log the error and return an error response.
		log.Debug("Error reading body: ", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Switch on request method
	switch req.Method {
	case http.MethodPost:
		log.Debug("/chatgroup POST")
		chatgroupNew(w, body)
	case http.MethodGet:
		log.Debug("/chatgroup GET")
		chatgroupGet(w, body)
	}
}

func messageNew(w http.ResponseWriter, body []byte) {
	// Decode body into dict
	var message_dict map[string]interface{}
	err := json.Unmarshal(body, &message_dict)
	if err != nil { // Error unmarshalling body
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Check user credentials
	username := message_dict["username"].(string)
	password := message_dict["password"].(string)

	// Check if user exists and password is correct
	rows, err := DB.Query("SELECT * FROM User WHERE username = ?;", username)
	if err != nil { // Error querying database
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	if !rows.Next() { // User does not exist
		log.Debug("Error: User does not exist")
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	var DB_user lib.User
	var joined_raw string
	err = rows.Scan( // Scan DB User
		&DB_user.Username,
		&DB_user.Firstname,
		&DB_user.Lastname,
		&DB_user.Email,
		&DB_user.Password,
		&joined_raw,
	)
	if err != nil { // Error scanning row
		log.Error("Error scanning row: ", err)
		http.Error(w, "can't scan row", http.StatusInternalServerError)
		return
	}

	if DB_user.Password != password { // Wrong password
		log.Debug("Error: Wrong password")
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	// prase message out of dict
	group := message_dict["group"].(string)
	message := message_dict["message"].(string)
	SentOn := time.Now()

	// Get all Group members
	rows, err = DB.Query(
		"SELECT user_name FROM ChatGroupMemberShip WHERE `chat_group_name` = ?;",
		group,
	)
	if err != nil { // Error querying database
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	// Iterate over all rows
	for rows.Next() {
		// Get user_name from row
		var user_name string
		err = rows.Scan(&user_name)
		if err != nil { // Error scanning row
			log.Error("Error scanning row: ", err)
			http.Error(w, "can't scan row", http.StatusInternalServerError)
			return
		}

		// Insert message into database
		stmt, err := DB.Prepare(
			"INSERT INTO Message (from_user_name, for_user_name, in_chat_group_name, message, sent_on) VALUES (?, ?, ?, ?, ?);",
		)
		if err != nil { // Error preparing statement
			log.Error("Error preparing statement: ", err)
			http.Error(w, "can't prepare statement", http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		// Execute statement
		_, err = stmt.Exec(username, user_name, group, message, SentOn)
		if err != nil { // Error executing statement
			log.Error("Error executing statement: ", err)
			http.Error(w, "can't execute statement", http.StatusInternalServerError)
			return
		}

	}
}

func messagesGet(w http.ResponseWriter, body []byte) {
	// decode body into user
	var request_user lib.User
	err := json.Unmarshal(body, &request_user)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Check if user exists and password is correct
	rows, err := DB.Query("SELECT * FROM User WHERE username = ?;", request_user.Username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}
	if !rows.Next() {
		log.Debug("Error: User does not exist")
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	var user lib.User
	var joined_raw string
	err = rows.Scan( // Readout DB User
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&joined_raw,
	)
	if err != nil { // Error scanning row
		log.Error("Error scanning row: ", err)
		http.Error(w, "can't scan row", http.StatusInternalServerError)
		return
	}

	if user.Password != request_user.Password { // Wrong password
		log.Debug("Error: Wrong password")
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	// Get Messages
	rows, err = DB.Query("SELECT * FROM Message WHERE `for_user_name` = ?;", request_user.Username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var messages []lib.Message

	// Iterate over all rows
	for rows.Next() {
		var message lib.Message
		var sent_on_raw string

		err = rows.Scan(&message.From, &message.To, &message.Group, &message.Message, &sent_on_raw)
		if err != nil { // Error scanning row
			log.Error("Error scanning row: ", err)
			http.Error(w, "can't scan row", http.StatusInternalServerError)
			return
		}

		// Parse sent_on_raw into time.Time
		message.SentOn, err = time.Parse(time.DateTime, sent_on_raw)
		if err != nil { // Error parsing database date
			log.Error("Error parsing database date: ", err)
			http.Error(w, "can't parse database date", http.StatusInternalServerError)
			return
		}

		messages = append(messages, message)
	}

	// Once retrieved, delete messages from database
	stmt, err := DB.Prepare("DELETE FROM Message WHERE `for_user_name` = ?;")
	if err != nil { // Error preparing statement
		log.Error("Error preparing statement: ", err)
		http.Error(w, "can't prepare statement", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(request_user.Username)
	if err != nil { // Error executing statement
		log.Error("Error executing statement: ", err)
		http.Error(w, "can't execute statement", http.StatusInternalServerError)
		return
	}

	// Return messages as json
	messages_json, err := json.Marshal(messages)
	if err != nil {
		log.Error("Error marshalling messages: ", err)
		http.Error(w, "can't marshal messages", http.StatusInternalServerError)
		return
	}

	// write response
	w.Write(messages_json)
	w.Header().Set("Content-Type", "application/json")
}

func message_handler(w http.ResponseWriter, req *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Debug("Error reading body: ", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	switch req.Method { // route based on request method
	case "POST":
		log.Debug("/message POST")
		messageNew(w, body)
	case "GET":
		log.Debug("/message GET")
		messagesGet(w, body)
	}
}

// TODO: Add authentication to all requests

func auth_handler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Debug("Error reading body: ", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Decode body into user
	var user lib.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Debug("Error unmarshalling body: ", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}

	// Get username and password from user
	username := user.Username
	password := user.Password

	// Validate username and password
	if username == "" || password == "" {
		log.Debug("Error: Username or password is empty")
		http.Error(w, "username or password is empty", http.StatusBadRequest)
		return
	}

	// Query database for user
	rows, err := DB.Query("SELECT * FROM User WHERE username = ?;", username)
	if err != nil {
		log.Error("Error querying database: ", err)
		http.Error(w, "can't query database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Check if user exists
	if !rows.Next() {
		log.Debug("Error: User not found")
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	// Scan user data from database
	var joined_raw string
	err = rows.Scan(
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&joined_raw,
	)

	if err != nil {
		log.Error("Error scanning row: ", err)
		http.Error(w, "can't scan row", http.StatusInternalServerError)
		return
	}

	// Parse joined date from database
	joined, err := time.Parse(time.DateTime, joined_raw)
	if err != nil {
		log.Error("Error parsing database date: ", err)
		http.Error(w, "can't parse database date", http.StatusInternalServerError)
		return
	}

	// Set joined date on user
	user.Joined = joined

	// Check if password is correct
	if user.Password != password {
		log.Debug("Error: Wrong password")
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	// Send okay
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.SetLevel(log.DebugLevel)
	// Connect to database GetDB()
	log.Debug("Connecting to database...")
	GetDB()

	defer DB.Close() // Close database connection when main function returns

	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/chatgroup", chatgroup_handler)
	http.HandleFunc("/message", message_handler)
	http.HandleFunc("/auth", auth_handler)

	http.ListenAndServe("localhost:8080", nil)

}
