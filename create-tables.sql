DROP TABLE IF EXISTS User;
CREATE TABLE User (
    username VARCHAR(10) NOT NULL,
    first_name VARCHAR(50) ,
    last_name VARCHAR(50) ,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    joined_on DATETIME NOT NULL,
    PRIMARY KEY (username)
);

-- Create admin user
INSERT INTO User 
    (username, first_name, last_name, email, password, joined_on) 
VALUES 
    ('admin', 'admin', 'admin', 'admin@localhost', 'admin', NOW());

DROP TABLE IF EXISTS ChatGroup;
CREATE TABLE ChatGroup (
    name VARCHAR(50) NOT NULL,
    created_on DATETIME NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    FOREIGN KEY (created_by) REFERENCES User(username),
    PRIMARY KEY (name)
);

-- Create General group
INSERT INTO ChatGroup (name, created_on, created_by) VALUES ('General', NOW(), 'admin');

DROP TABLE IF EXISTS ChatGroupMemberShip;
CREATE TABLE ChatGroupMemberShip (
    user_name VARCHAR(10) NOT NULL,
    chat_group_name VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_name) REFERENCES User(username),
    FOREIGN KEY (chat_group_name) REFERENCES ChatGroup(name),
    PRIMARY KEY (user_name, chat_group_name)
);

INSERT INTO ChatGroupMemberShip (user_name, chat_group_name) VALUES ('admin', 'General');

-- admin is a member of General group

DROP TABLE IF EXISTS Message;
CREATE TABLE Message (
    from_user_name VARCHAR(10) NOT NULL,
    for_user_name VARCHAR(10) NOT NULL,
    in_chat_group_name VARCHAR(10) NOT NULL,
    message VARCHAR(255) NOT NULL,
    sent_on DATETIME NOT NULL,
    PRIMARY KEY (from_user_name,for_user_name, in_chat_group_name, sent_on),
    FOREIGN KEY (from_user_name) REFERENCES User(username),
    FOREIGN KEY (for_user_name) REFERENCES User(username),
    FOREIGN KEY (in_chat_group_name) REFERENCES ChatGroup(name)
);

INSERT INTO Message (from_user_name,for_user_name,in_chat_group_name,message,sent_on) VALUES ('admin','admin','General','Welcome to Chat App',NOW());

