CREATE TABLE IF NOT EXISTS 
auth (
    username TEXT PRIMARY KEY, 
    password TEXT, 
    email TEXT
);

CREATE TABLE IF NOT EXISTS 
auth_device (
    device INTEGER PRIMARY KEY AUTOINCREMENT,
    auth TEXT,
    platform TEXT,
    token TEXT, 
    FOREIGN KEY(auth) REFERENCES auth(username)
);

