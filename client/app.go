package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

type ChatGroup struct {
	Name      string    `json:"name,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	Members   []string  `json:"members,omitempty"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	From   string    `json:"from,omitempty"`
	To     string    `json:"to,omitempty"`
	Text   string    `json:"text,omitempty"`
	SentOn time.Time `json:"sent_on,omitempty"`
}

// DEV CONSTANTS
const (
	SERVER_IP   = "localhost"
	SERVER_PORT = "8080"
	LOCAL_DB    = "/Users/jakobsachs/.cache/chatservice/db.json"
	USERNAME    = "admin"
	PASSWORD    = "admin"
)

// App struct
type App struct {
	ctx        context.Context
	logger     logger.Logger
	db_path    string
	ChatGroups []ChatGroup
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Read the database file and return the of ChatGroup
func read_db(app *App) []ChatGroup {
	// Read the file
	data, err := os.ReadFile(LOCAL_DB)
	if err != nil {
		panic(err)
	}

	// Log the data
	app.logger.Debug("Read the following data from the database:")
	app.logger.Debug(string(data))

	// read content to list of groups
	var groups []ChatGroup
	err = json.Unmarshal(data, &groups)
	if err != nil {
		panic(err)
	}

	return groups
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// Save the context
	a.ctx = ctx

	// Check if the database file exists
	if _, err := os.Stat(LOCAL_DB); os.IsNotExist(err) {
		// Create the file
		f, err := os.Create(LOCAL_DB)
		if err != nil {
			panic(err)
		}

		// Write empty array to the file
		os.WriteFile(LOCAL_DB, []byte("[]"), 0644)
		f.Close()
	}

	// Read the database file
	a.ChatGroups = read_db(a)
	println(a.ChatGroups)
}
