package lib

import "time"

type User struct {
	Username  string    `json:"username,omitempty"`
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Joined    time.Time `json:"joined,omitempty"`
}

type ChatGroup struct {
	Name      string    `json:"name,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	Members   []string  `json:"members,omitempty"`
}

type Message struct {
	From    string    `json:"from,omitempty"`
	To      string    `json:"to,omitempty"`
	Group   string    `json:"group,omitempty"`
	Message string    `json:"message,omitempty"`
	SentOn  time.Time `json:"sent_on,omitempty"`
}
