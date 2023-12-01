# go_chatservice
A Mostly abandoned project to attempt to write a chat service & Desktop client mostly in Golang.

 - **Server**: The server code is contained in the `server.go` file, which should probably be 
 split up into submodules. 
 - **Client**: Client code is located in the `client/` directory, and is its own project.

 The server runs and handles incomming connections for chatgroups and messages. 
 It manages to a my-sql DB to save & retrieve the data.

