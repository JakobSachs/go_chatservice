package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
)

// "github.com/charmbracelet/bubbles/table"
// tea "github.com/charmbracelet/bubbletea"
// "github.com/charmbracelet/lipgloss"

// static server settings
const (
	SERVER_IP   = "localhost"
	SERVER_PORT = "8080"
	LOCAL_DB    = "db.json"
)

// state of the app
type app_state struct {
	db_file os.File
}

var state app_state // global state of the app

func auth(username string, password string) error {
	// Create JSON for request
	json := []byte(`{"username": "` + username + `", "password": "` + password + `"}`)
	reader := bytes.NewReader(json)

	// Connect to the server to validate the credentials
	req, err := http.NewRequest("GET", "http://"+SERVER_IP+":"+SERVER_PORT+"/auth", reader)
	if err != nil {
		fmt.Printf("Error connecting to server %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to server %v\n", err)
		return err
	}

	// Check if the credentials are valid
	if res.StatusCode != 200 {
		fmt.Println("Invalid credentials")
		return errors.New("Invalid credentials")
	}

	return nil
}

// Setup the local json database for the client
// If the database/file does not exist, create it
func setupLocalDB() (os.File, error) {
	// Check if the file exists
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return os.File{}, err
	}

	// Check if the file exists
	if _, err := os.Stat(home_dir + "/" + LOCAL_DB); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create(home_dir + "/" + LOCAL_DB)
		if err != nil {
			return os.File{}, err
		}

		// Close the file
		defer file.Close()
	}

	// Open the file
	file, err := os.OpenFile(home_dir+"/"+LOCAL_DB, os.O_RDWR, 0644)
	if err != nil {
		return os.File{}, err
	}

	return *file, nil
}

// Read the local json database for the client and return the data
func getLocalDBContent() (map[string]string, error) {
	// Read file content
	file_content, err := ioutil.ReadAll(&state.db_file)
	if err != nil {
		return nil, err
	}

	// Convert the file content to a map
	var data map[string]string
	err = json.Unmarshal(file_content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	// Setup the local database
	file, err := setupLocalDB()
	if err != nil {
		fmt.Printf("Error setting up local database %v\n", err)
		return
	}
	state.db_file = file

	// Prompt for the username
	username_prompt := promptui.Prompt{
		Label:    "Username",
		Validate: nil,
	}

	username, err := username_prompt.Run() // Get the username
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Prompt for the password
	password_prompt := promptui.Prompt{
		Label:    "Password",
		Validate: nil,
		Mask:     '*',
	}

	password, err := password_prompt.Run() // Get the password
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Authenticate the user with given credentials
	err = auth(password, username)
	if err == nil {
		fmt.Println("Login successful ✅")
	} else {
		fmt.Println("Login failed ❌")
	}
}
